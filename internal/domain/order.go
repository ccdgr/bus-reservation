package domain

import (
	"context"
	"time"
)

// 订单状态常量
const (
	StatusPendingPayment      = 0 // 待支付
	StatusPendingVerification = 1 // 待核验 (已支付)
	StatusCancelled           = 2 // 已取消
	StatusExpired             = 3 // 已过期
	StatusVerified            = 4 // 已核验
	StatusRefunding           = 5 // 退款中
)

// Order 代表预约订单
type Order struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	BusID     uint64    `json:"bus_id" gorm:"index"`
	PaymentID string    `json:"payment_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联字段 (GORM 自动加载)
	Bus *Bus `json:"bus,omitempty" gorm:"foreignKey:BusID"`
}

// OrderRepository 订单数据访问接口
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uint64) (*Order, error)
	UpdateStatus(ctx context.Context, orderID uint64, status int) error
	UpdateStatusAndPaymentID(ctx context.Context, orderID uint64, status int, paymentID string) error
	ListByUserID(ctx context.Context, userID uint64) ([]*Order, error)
	CheckUserHasActiveOrder(ctx context.Context, userID, busID uint64) (bool, error)
}

// OrderUsecase 订单业务逻辑接口
type OrderUsecase interface {
	Create(ctx context.Context, userID, busID uint64) (*Order, error)
	ListByUserID(ctx context.Context, userID uint64) ([]*Order, error)
	Cancel(ctx context.Context, orderID uint64) error
	Pay(ctx context.Context, orderID uint64) (string, error)
	CapturePayPalPayment(ctx context.Context, orderID uint64, paypalToken string) error
	Verify(ctx context.Context, orderID uint64) error
}

// OrderStateMachine 状态机逻辑
func CanTransition(from, to int) bool {
	switch from {
	case StatusPendingPayment:
		// 待支付可以去：待核验(支付)、已取消(手动/超时)、已过期(发车)
		return to == StatusPendingVerification || to == StatusCancelled || to == StatusExpired
	case StatusPendingVerification:
		// 待核验可以去：已核验(乘车)、退款中(退票)、已过期(漏乘)
		return to == StatusVerified || to == StatusRefunding || to == StatusExpired
	case StatusRefunding:
		// 退款中可以去：已取消(退款成功)
		return to == StatusCancelled
	default:
		// 其他终态 (Cancelled, Expired, Verified) 不允许再流转
		return false
	}
}

func GetStatusName(status int) string {
	switch status {
	case StatusPendingPayment:
		return "待支付"
	case StatusPendingVerification:
		return "待核验"
	case StatusCancelled:
		return "已取消"
	case StatusExpired:
		return "已过期"
	case StatusVerified:
		return "已核验"
	case StatusRefunding:
		return "退款中"
	default:
		return "未知"
	}
}

// CheckAndFixStatus 根据发车时间自动修正状态 (逻辑上的“已过期”)
func (o *Order) CheckAndFixStatus(departureTime time.Time) bool {
	now := time.Now()
	if now.After(departureTime) {
		if o.Status == StatusPendingPayment || o.Status == StatusPendingVerification {
			o.Status = StatusExpired
			return true
		}
	}
	return false
}
