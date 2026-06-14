package mq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/ccdgr/bus-reservation/pkg/payment"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderConsumer struct {
	conn         *amqp.Connection
	orderRepo    domain.OrderRepository
	busRepo      domain.BusRepository
	redisRepo    domain.BusRepository
	paypalClient *payment.PayPalClient
}

func NewOrderConsumer(conn *amqp.Connection, orderRepo domain.OrderRepository, busRepo domain.BusRepository, redisRepo domain.BusRepository, paypalClient *payment.PayPalClient) *OrderConsumer {
	return &OrderConsumer{
		conn:         conn,
		orderRepo:    orderRepo,
		busRepo:      busRepo,
		redisRepo:    redisRepo,
		paypalClient: paypalClient,
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
	_, err = ch.QueueDeclare("order_delay", true, false, false, false, args)
	if err != nil {
		return err
	}

	cancelMsgs, err := ch.Consume(qCancel.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// 3. 声明退款队列
	qRefund, err := ch.QueueDeclare("order_refund", true, false, false, false, nil)
	if err != nil {
		return err
	}

	refundMsgs, err := ch.Consume(qRefund.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// 协程 A: 处理退款
	go func() {
		for d := range refundMsgs {
			var msg OrderMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				d.Ack(false)
				continue
			}

			order, err := c.orderRepo.GetByID(ctx, msg.OrderID)
			if err != nil {
				slog.Error("处理退款时无法获取订单", "order_id", msg.OrderID, "error", err)
				d.Nack(false, true)
				continue
			}

			if order.Status == domain.StatusRefunding {
				if order.PaymentID != "" && c.paypalClient != nil {
					err = c.paypalClient.RefundOrder(ctx, order.PaymentID)
					if err != nil {
						slog.Error("PayPal 退款接口调用失败", "order_id", order.ID, "payment_id", order.PaymentID, "error", err)
						d.Nack(false, true) // 退款失败重回队列重试
						continue
					}
				}

				// 即使没有 payment_id (mock payment)，我们也将其标记为退款成功并恢复库存
				err = c.orderRepo.UpdateStatus(ctx, order.ID, domain.StatusCancelled)
				if err != nil {
					slog.Error("更新订单为已取消状态失败", "order_id", order.ID, "error", err)
					d.Nack(false, true)
					continue
				}

				c.redisRepo.IncrSeat(ctx, order.BusID)
				c.busRepo.UpdateSeat(ctx, order.BusID, 1)

				slog.Info("订单退款并取消成功", "order_id", order.ID)
			} else {
				slog.Info("订单状态并非退款中，跳过", "order_id", order.ID)
			}

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
