package usecase
import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ccdgr/bus-reservation/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/smartwalle/alipay/v3"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
	busRepo   domain.BusRepository
	redisRepo domain.BusRepository
	mqChannel *amqp.Channel
	aliClient *alipay.Client
	notifyURL string
	returnURL string
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	busRepo domain.BusRepository,
	redisRepo domain.BusRepository,
	mqChannel *amqp.Channel,
	aliClient *alipay.Client,
	notifyURL, returnURL string,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
		busRepo:   busRepo,
		redisRepo: redisRepo,
		mqChannel: mqChannel,
		aliClient: aliClient,
		notifyURL: notifyURL,
		returnURL: returnURL,
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

	msg := OrderMessage{
		UserID: userID,
		BusID:  busID,
	}
	body, _ := json.Marshal(msg)

	err = u.mqChannel.PublishWithContext(ctx,
		"",
		"order_jobs",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		slog.Error("failed to publish order message", "error", err)
		u.redisRepo.IncrSeat(ctx, busID)
		return nil, err
	}

	return &domain.Order{
		UserID: userID,
		BusID:  busID,
		Status: domain.StatusPendingPayment,
		Bus:    bus,
	}, nil
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

	if u.aliClient == nil {
		// Mock payment logic if alipay is not configured
		err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusPendingVerification)
		return "", err
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = u.notifyURL
	p.ReturnURL = u.returnURL
	p.Subject = fmt.Sprintf("Bus Reservation Order #%d", orderID)
	p.OutTradeNo = strconv.FormatUint(orderID, 10)
	p.TotalAmount = "0.01" // Simulated amount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := u.aliClient.TradePagePay(p)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (u *orderUsecase) ConfirmPayment(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		// Might already be verified/paid if callback arrives multiple times
		if order.Status == domain.StatusPendingVerification {
			return nil
		}
		return fmt.Errorf("current status %s cannot confirm payment", domain.GetStatusName(order.Status))
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
