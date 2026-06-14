package usecase
import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/ccdgr/bus-reservation/pkg/payment"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/smartwalle/alipay/v3"
)

type orderUsecase struct {
	orderRepo    domain.OrderRepository
	busRepo      domain.BusRepository
	redisRepo    domain.BusRepository
	mqChannel    *amqp.Channel
	paypalClient *payment.PayPalClient
	aliClient    *alipay.Client
	notifyURL    string
	returnURL    string
	cancelURL    string
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	busRepo domain.BusRepository,
	redisRepo domain.BusRepository,
	mqChannel *amqp.Channel,
	paypalClient *payment.PayPalClient,
	aliClient *alipay.Client,
	notifyURL, returnURL, cancelURL string,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:    orderRepo,
		busRepo:      busRepo,
		redisRepo:    redisRepo,
		mqChannel:    mqChannel,
		paypalClient: paypalClient,
		aliClient:    aliClient,
		notifyURL:    notifyURL,
		returnURL:    returnURL,
		cancelURL:    cancelURL,
	}
}


type OrderMessage struct {
	UserID  uint64 `json:"user_id"`
	BusID   uint64 `json:"bus_id"`
	OrderID uint64 `json:"order_id,omitempty"`
}

func (u *orderUsecase) Create(ctx context.Context, userID, busID uint64) (*domain.Order, error) {
	// 1. 业务规则校验：同一用户同一车次只允许购买一次（限购）
	hasActiveOrder, err := u.orderRepo.CheckUserHasActiveOrder(ctx, userID, busID)
	if err != nil {
		slog.Error("检查用户活跃订单失败", "user_id", userID, "bus_id", busID, "error", err)
		return nil, err
	}
	if hasActiveOrder {
		return nil, fmt.Errorf("您已经预定了该车次，请勿重复购买")
	}

	bus, err := u.busRepo.GetByID(ctx, busID)
	if err != nil {
		return nil, fmt.Errorf("未找到该班次")
	}

	ok, err := u.redisRepo.DecrSeat(ctx, busID)
	if err != nil {
		slog.Error("Redis 扣减库存失败", "bus_id", busID, "error", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("该班次已售罄")
	}

	// 同步创建订单以获取 ID
	order := &domain.Order{
		UserID: userID,
		BusID:  busID,
		Status: domain.StatusPendingPayment,
	}
	if err := u.orderRepo.Create(ctx, order); err != nil {
		slog.Error("在数据库中创建订单失败", "error", err)
		u.redisRepo.IncrSeat(ctx, busID)
		return nil, err
	}

	// 异步更新物理库存
	go func() {
		if err := u.busRepo.UpdateSeat(context.Background(), busID, -1); err != nil {
			slog.Error("更新数据库物理库存失败", "error", err)
		}
	}()

	// 发送延时取消消息
	delayMsg := OrderMessage{
		OrderID: order.ID,
		BusID:   order.BusID,
	}
	delayBody, _ := json.Marshal(delayMsg)
	err = u.mqChannel.PublishWithContext(ctx,
		"",
		"order_delay",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         delayBody,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		slog.Error("发送延时取消消息失败", "error", err)
		// 不直接失败，因为订单已创建，只是可能无法自动超时取消
	}

	order.Bus = bus
	return order, nil
}

func (u *orderUsecase) ListByUserID(ctx context.Context, userID uint64) ([]*domain.Order, error) {
	orders, err := u.orderRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, o := range orders {
		if o.Bus != nil {
			o.CheckAndFixStatus(o.Bus.StartTime)
		}
	}

	return orders, nil
}

func (u *orderUsecase) Cancel(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == domain.StatusPendingVerification {
		// 已支付状态取消，进入退款中状态
		if !domain.CanTransition(order.Status, domain.StatusRefunding) {
			return fmt.Errorf("当前状态 [%s] 无法执行退票操作", domain.GetStatusName(order.Status))
		}

		err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusRefunding)
		if err != nil {
			return err
		}

		// 发送退款消息到 MQ
		refundMsg := OrderMessage{
			OrderID: order.ID,
			BusID:   order.BusID,
		}
		refundBody, _ := json.Marshal(refundMsg)
		err = u.mqChannel.PublishWithContext(ctx,
			"",
			"order_refund",
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         refundBody,
				DeliveryMode: amqp.Persistent,
			})
		if err != nil {
			slog.Error("发送退款消息失败", "order_id", orderID, "error", err)
			return fmt.Errorf("提交退款申请失败，请稍后再试")
		}
		return nil
	}

	if !domain.CanTransition(order.Status, domain.StatusCancelled) {
		return fmt.Errorf("当前状态 [%s] 无法执行取消操作", domain.GetStatusName(order.Status))
	}

	err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusCancelled)
	if err != nil {
		return err
	}

	u.redisRepo.IncrSeat(ctx, order.BusID)
	return u.busRepo.UpdateSeat(ctx, order.BusID, 1)
}

func (u *orderUsecase) Pay(ctx context.Context, orderID uint64, paymentMethod string) (string, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return "", err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		return "", fmt.Errorf("当前状态 [%s] 无法发起支付", domain.GetStatusName(order.Status))
	}

	if paymentMethod == "paypal" && u.paypalClient != nil {
		params := payment.CreateOrderParams{
			Amount:    "5.00",
			ReturnURL: fmt.Sprintf("%s?order_id=%d", u.returnURL, orderID),
			CancelURL: u.cancelURL,
		}
		url, err := u.paypalClient.CreateOrder(ctx, params)
		if err != nil {
			slog.Error("调用 PayPal 创建订单失败", "error", err)
			return "", err
		}
		return url, nil
	} else if paymentMethod == "alipay" && u.aliClient != nil {
		var p = alipay.TradePagePay{}
		p.NotifyURL = u.notifyURL
		p.ReturnURL = u.cancelURL // frontend orders page
		p.Subject = fmt.Sprintf("Bus Reservation Order #%d", orderID)
		p.OutTradeNo = strconv.FormatUint(orderID, 10)
		p.TotalAmount = "5.00" // Fixed 5 RMB amount
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"

		url, err := u.aliClient.TradePagePay(p)
		if err != nil {
			slog.Error("调用 Alipay 创建订单失败", "error", err)
			return "", err
		}
		return url.String(), nil
	}

	// Mock payment logic
	err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPendingVerification)
	return "", err
}

func (u *orderUsecase) ConfirmAlipayPayment(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		// Idempotent success
		if order.Status == domain.StatusPendingVerification {
			return nil
		}
		return fmt.Errorf("当前状态 [%s] 无法确认支付", domain.GetStatusName(order.Status))
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPendingVerification)
}

func (u *orderUsecase) CapturePayPalPayment(ctx context.Context, orderID uint64, paypalToken string) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		// Idempotent success
		if order.Status == domain.StatusPendingVerification {
			return nil
		}
		return fmt.Errorf("当前状态 [%s] 无法确认支付", domain.GetStatusName(order.Status))
	}

	// Call PayPal to capture the order
	captureID, err := u.paypalClient.CaptureOrder(ctx, paypalToken)
	if err != nil {
		slog.Error("PayPal 扣款确认失败", "error", err)
		return err
	}

	return u.orderRepo.UpdateStatusAndPaymentID(ctx, orderID, domain.StatusPendingVerification, captureID)
}

func (u *orderUsecase) Verify(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusVerified) {
		return fmt.Errorf("当前状态 [%s] 无法执行核验操作", domain.GetStatusName(order.Status))
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusVerified)
}
