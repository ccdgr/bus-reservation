package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ccdgr/bus-reservation/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
	busRepo   domain.BusRepository
	redisRepo domain.BusRepository
	mqChannel *amqp.Channel
}

func NewOrderUsecase(orderRepo domain.OrderRepository, busRepo domain.BusRepository, redisRepo domain.BusRepository, mqChannel *amqp.Channel) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
		busRepo:   busRepo,
		redisRepo: redisRepo,
		mqChannel: mqChannel,
	}
}

type OrderMessage struct {
	UserID uint64 `json:"user_id"`
	BusID  uint64 `json:"bus_id"`
}

func (u *orderUsecase) Create(ctx context.Context, userID, busID uint64) (*domain.Order, error) {
	// 校验班次是否存在
	bus, err := u.busRepo.GetByID(ctx, busID)
	if err != nil {
		return nil, fmt.Errorf("bus not found")
	}

	// 1. 原子扣减 Redis 库存
	ok, err := u.redisRepo.DecrSeat(ctx, busID)
	if err != nil {
		slog.Error("failed to decrease redis stock", "bus_id", busID, "error", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("no seats available")
	}

	// 2. 发送异步消息到 MQ
	msg := OrderMessage{
		UserID: userID,
		BusID:  busID,
	}
	body, _ := json.Marshal(msg)

	err = u.mqChannel.PublishWithContext(ctx,
		"",           // exchange
		"order_jobs", // routing key (queue name)
		true,         // mandatory
		false,        // immediate
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

	// 逻辑上修正过期状态（如果还未支付或待核验但发车时间已过）
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

	// 更新状态
	err = u.orderRepo.UpdateStatus(ctx, orderID, domain.StatusCancelled)
	if err != nil {
		return err
	}

	// 归还 Redis 和 DB 库存
	u.redisRepo.IncrSeat(ctx, order.BusID)
	return u.busRepo.UpdateSeat(ctx, order.BusID, 1)
}

func (u *orderUsecase) Pay(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(order.Status, domain.StatusPendingVerification) {
		return fmt.Errorf("current status %s cannot be paid", domain.GetStatusName(order.Status))
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
