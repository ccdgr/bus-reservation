package usecase
import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/ccdgr/bus-reservation/pkg/payment"
	amqp "github.com/rabbitmq/amqp091-go"
)

type orderUsecase struct {
	orderRepo    domain.OrderRepository
	busRepo      domain.BusRepository
	redisRepo    domain.BusRepository
	mqChannel    *amqp.Channel
	paypalClient *payment.PayPalClient
	returnURL    string
	cancelURL    string
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	busRepo domain.BusRepository,
	redisRepo domain.BusRepository,
	mqChannel *amqp.Channel,
	paypalClient *payment.PayPalClient,
	returnURL, cancelURL string,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:    orderRepo,
		busRepo:      busRepo,
		redisRepo:    redisRepo,
		mqChannel:    mqChannel,
		paypalClient: paypalClient,
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
	bus, err := u.busRepo.GetByID(ctx, busID)
	if err != nil {
		return nil, fmt.Errorf("bus not found")
	}

	ok, err := u.redisRepo.DecrSeat(ctx, busID)
	if err != nil {
		slog.Error("failed to decrease redis stock", "bus_id", busID, "error", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("no seats available")
	}

	// 同步创建订单以获取 ID
	order := &domain.Order{
		UserID: userID,
		BusID:  busID,
		Status: domain.StatusPendingPayment,
	}
	if err := u.orderRepo.Create(ctx, order); err != nil {
		slog.Error("failed to create order in db", "error", err)
		u.redisRepo.IncrSeat(ctx, busID)
		return nil, err
	}

	// 异步更新物理库存
	go func() {
		if err := u.busRepo.UpdateSeat(context.Background(), busID, -1); err != nil {
			slog.Error("failed to update bus seat in db", "error", err)
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
		slog.Error("failed to publish delay cancel message", "error", err)
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

	if !domain.CanTransition(order.Status, domain.StatusCancelled) {
		return fmt.Errorf("current status %s cannot be cancelled", domain.GetStatusName(order.Status))
	}

	err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusCancelled)
	if err != nil {
		return err
	}

	u.redisRepo.IncrSeat(ctx, order.BusID)
	return u.busRepo.UpdateSeat(ctx, order.BusID, 1)
}

func (u *orderUsecase) Pay(ctx context.Context, orderID uint64) (string, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return "", err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		return "", fmt.Errorf("current status %s cannot be paid", domain.GetStatusName(order.Status))
	}

	if u.paypalClient != nil {
		params := payment.CreateOrderParams{
			Amount:    "5.00",
			ReturnURL: fmt.Sprintf("%s?order_id=%d", u.returnURL, orderID),
			CancelURL: u.cancelURL,
		}
		url, err := u.paypalClient.CreateOrder(ctx, params)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	// Mock payment logic
	err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPendingVerification)
	return "", err
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
		return fmt.Errorf("current status %s cannot be captured", domain.GetStatusName(order.Status))
	}

	// Call PayPal to capture the order
	err = u.paypalClient.CaptureOrder(ctx, paypalToken)
	if err != nil {
		return err
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPendingVerification)
}

func (u *orderUsecase) Verify(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusVerified) {
		return fmt.Errorf("current status %s cannot be verified", domain.GetStatusName(order.Status))
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusVerified)
}
