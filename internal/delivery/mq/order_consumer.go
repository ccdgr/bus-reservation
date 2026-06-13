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
	redisRepo domain.BusRepository
}

func NewOrderConsumer(conn *amqp.Connection, orderRepo domain.OrderRepository, busRepo domain.BusRepository, redisRepo domain.BusRepository) *OrderConsumer {
	return &OrderConsumer{
		conn:      conn,
		orderRepo: orderRepo,
		busRepo:   busRepo,
		redisRepo: redisRepo,
	}
}

type OrderMessage struct {
	UserID  uint64 `json:"user_id"`
	BusID   uint64 `json:"bus_id"`
	OrderID uint64 `json:"order_id,omitempty"`
}

func (c *OrderConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 1. 声明延时交换机和队列逻辑
	err = ch.ExchangeDeclare("order_dlx", "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	qCancel, err := ch.QueueDeclare("order_cancel", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(qCancel.Name, "cancel", "order_dlx", false, nil)

	args := amqp.Table{
		"x-dead-letter-exchange":    "order_dlx",
		"x-dead-letter-routing-key": "cancel",
		"x-message-ttl":             int32(15 * 60 * 1000), // 15 分钟
	}
	qDelay, err := ch.QueueDeclare("order_delay", true, false, false, false, args)
	if err != nil {
		return err
	}

	// 2. 声明普通入库队列
	qJobs, err := ch.QueueDeclare("order_jobs", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(qJobs.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	cancelMsgs, err := ch.Consume(qCancel.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// 协程 A: 处理订单入库并发送延时取消消息
	go func() {
		for d := range msgs {
			var msg OrderMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				d.Nack(false, false)
				continue
			}

			order := &domain.Order{
				UserID: msg.UserID,
				BusID:  msg.BusID,
				Status: domain.StatusPendingPayment,
			}
			if err := c.orderRepo.Create(ctx, order); err != nil {
				slog.Error("failed to create order in db", "error", err)
				d.Nack(false, true)
				continue
			}

			if err := c.busRepo.UpdateSeat(ctx, msg.BusID, -1); err != nil {
				slog.Error("failed to update bus seat in db", "error", err)
			}

			delayMsg := OrderMessage{
				OrderID: order.ID,
				BusID:   order.BusID,
			}
			delayBody, _ := json.Marshal(delayMsg)
			ch.PublishWithContext(ctx, "", qDelay.Name, false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         delayBody,
				DeliveryMode: amqp.Persistent,
			})

			slog.Info("order created and delay task scheduled", "order_id", order.ID)
			d.Ack(false)
		}
	}()

	// 协程 B: 处理过期订单取消
	go func() {
		for d := range cancelMsgs {
			var msg OrderMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				d.Ack(false)
				continue
			}

			order, err := c.orderRepo.GetByID(ctx, msg.OrderID)
			if err != nil {
				slog.Error("failed to fetch order for delay cancel", "order_id", msg.OrderID, "error", err)
				d.Nack(false, true)
				continue
			}

			if order.Status == domain.StatusPendingPayment {
				err = c.orderRepo.UpdateStatus(ctx, order.ID, domain.StatusCancelled)
				if err != nil {
					slog.Error("failed to auto-cancel order", "order_id", order.ID, "error", err)
					d.Nack(false, true)
					continue
				}

				c.redisRepo.IncrSeat(ctx, order.BusID)
				c.busRepo.UpdateSeat(ctx, order.BusID, 1)
				
				slog.Info("order auto-cancelled due to timeout", "order_id", order.ID)
			} else {
				slog.Info("order already paid or cancelled, skipping auto-cancel", "order_id", order.ID)
			}
			
			d.Ack(false)
		}
	}()

	slog.Info("MQ consumers started (Jobs & Delay Cancel)")
	<-ctx.Done()
	return nil
}
