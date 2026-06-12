package domain

import (
	"context"
	"time"
)

// User 代表系统用户（学生或教职工）
type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Password  string    `json:"-"` // 密码不参与 JSON 序列化
	RealName  string    `json:"real_name"`
	UserType  int       `json:"user_type"` // 0: 学生, 1: 教职工
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}
