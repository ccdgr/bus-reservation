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

// BusRepository 班次数据访问接口
type BusRepository interface {
	GetByID(ctx context.Context, id uint64) (*Bus, error)
	List(ctx context.Context, origin, dest string, date string) ([]*Bus, error)
	UpdateSeat(ctx context.Context, busID uint64, delta int) error
	DecrSeat(ctx context.Context, busID uint64) (bool, error)
	IncrSeat(ctx context.Context, busID uint64) error
	GetStock(ctx context.Context, busID uint64) (int, error)
}

// BusUsecase 班次业务逻辑接口
type BusUsecase interface {
	List(ctx context.Context, origin, dest string, date string) ([]*Bus, error)
	GetByID(ctx context.Context, id uint64) (*Bus, error)
}
