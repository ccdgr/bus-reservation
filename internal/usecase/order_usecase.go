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
	redisRepo domain.BusRepository // Special implementation for DecrSeat
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
		// 回滚 Redis 库存
		u.redisRepo.IncrSeat(ctx, busID)
		return nil, err
	}

	// 3. 返回一个预备订单（实际入库由消费者完成，这里可以返回一个 ID 生成器生成的 ID 或者占位符）
	// 为简化，我们这里只返回一个带有关键信息的结构体，前端可以轮询查询结果
	return &domain.Order{
		UserID: userID,
		BusID:  busID,
		Status: 0, // Pending
	}, nil
}

func (u *orderUsecase) ListByUserID(ctx context.Context, userID uint64) ([]*domain.Order, error) {
	return u.orderRepo.ListByUserID(ctx, userID)
}

func (u *orderUsecase) Cancel(ctx context.Context, orderID uint64) error {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status != 0 {
		return fmt.Errorf("order cannot be cancelled")
	}

	// 更新状态为已取消
	err = u.orderRepo.UpdateStatus(ctx, orderID, 2)
	if err != nil {
		return err
	}

	// 归还库存
	return u.busRepo.UpdateSeat(ctx, order.BusID, 1)
}
