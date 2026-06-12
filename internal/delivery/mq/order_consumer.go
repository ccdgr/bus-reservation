package mq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ccdgr/bus-reservation/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderConsumer struct {
	conn      *amqp.Connection
	orderRepo domain.OrderRepository
	busRepo   domain.BusRepository
}

func NewOrderConsumer(conn *amqp.Connection, orderRepo domain.OrderRepository, busRepo domain.BusRepository) *OrderConsumer {
	return &OrderConsumer{
		conn:      conn,
		orderRepo: orderRepo,
		busRepo:   busRepo,
	}
}

type OrderMessage struct {
	UserID uint64 `json:"user_id"`
	BusID  uint64 `json:"bus_id"`
}

func (c *OrderConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_jobs", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg OrderMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				slog.Error("failed to unmarshal message", "error", err)
				d.Nack(false, false)
				continue
			}

			// 1. 数据库持久化订单
			order := &domain.Order{
				UserID: msg.UserID,
				BusID:  msg.BusID,
				Status: 1, // 直接设为已支付/成功，简化逻辑
			}
			if err := c.orderRepo.Create(ctx, order); err != nil {
				slog.Error("failed to create order in db", "error", err)
				d.Nack(false, true) // 重回队列
				continue
			}

			// 2. 更新数据库库存
			if err := c.busRepo.UpdateSeat(ctx, msg.BusID, -1); err != nil {
				slog.Error("failed to update bus seat in db", "error", err)
				// 这里可能需要补偿逻辑
			}

			slog.Info("order processed successfully", "user_id", msg.UserID, "bus_id", msg.BusID)
			d.Ack(false)
		}
	}()

	slog.Info("MQ consumer started")
	<-ctx.Done()
	return nil
}
