package domain

import (
	"context"
	"time"
)

// Bus 代表校车班次
type Bus struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Number    string    `json:"number" gorm:"uniqueIndex"` // 班次号
	Origin    string    `json:"origin"`                   // 起点
	Dest      string    `json:"dest"`                     // 终点
	StartTime time.Time `json:"start_time"`               // 发车时间
	TotalSeat int       `json:"total_seat"`               // 总座位数
	LeftSeat  int       `json:"left_seat"`                // 剩余座位数
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order 代表预约订单
type Order struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	BusID     uint64    `json:"bus_id" gorm:"index"`
	Status    int       `json:"status"` // 0: 待支付, 1: 已支付, 2: 已取消
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BusRepository 班次数据访问接口
type BusRepository interface {
	GetByID(ctx context.Context, id uint64) (*Bus, error)
	List(ctx context.Context) ([]*Bus, error)
	UpdateSeat(ctx context.Context, busID uint64, delta int) error
}

// OrderRepository 订单数据访问接口
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uint64) (*Order, error)
	UpdateStatus(ctx context.Context, orderID uint64, status int) error
}
