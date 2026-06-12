package domain

import (
	"context"
	"time"
)

// Order 代表预约订单
type Order struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	BusID     uint64    `json:"bus_id" gorm:"index"`
	Status    int       `json:"status"` // 0: 待支付, 1: 已支付, 2: 已取消
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderRepository 订单数据访问接口
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uint64) (*Order, error)
	UpdateStatus(ctx context.Context, orderID uint64, status int) error
	ListByUserID(ctx context.Context, userID uint64) ([]*Order, error)
}

// OrderUsecase 订单业务逻辑接口
type OrderUsecase interface {
	Create(ctx context.Context, userID, busID uint64) (*Order, error)
	ListByUserID(ctx context.Context, userID uint64) ([]*Order, error)
	Cancel(ctx context.Context, orderID uint64) error
}
