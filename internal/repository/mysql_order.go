package repository

import (
	"context"
	"github.com/ccdgr/bus-reservation/internal/domain"
	"gorm.io/gorm"
)

type mysqlOrderRepository struct {
	db *gorm.DB
}

func NewMySQLOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &mysqlOrderRepository{db: db}
}

func (r *mysqlOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *mysqlOrderRepository) GetByID(ctx context.Context, id uint64) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).Preload("Bus").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *mysqlOrderRepository) UpdateStatus(ctx context.Context, orderID uint64, status int) error {
	return r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *mysqlOrderRepository) ListByUserID(ctx context.Context, userID uint64) ([]*domain.Order, error) {
	var orders []*domain.Order
	if err := r.db.WithContext(ctx).Preload("Bus").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
