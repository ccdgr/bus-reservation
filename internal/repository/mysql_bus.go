package repository

import (
	"context"
	"fmt"
	"github.com/ccdgr/bus-reservation/internal/domain"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type mysqlBusRepository struct {
	db *gorm.DB
	sg singleflight.Group
}

func NewMySQLBusRepository(db *gorm.DB) domain.BusRepository {
	return &mysqlBusRepository{db: db}
}

func (r *mysqlBusRepository) GetByID(ctx context.Context, id uint64) (*domain.Bus, error) {
	key := fmt.Sprintf("bus:%d", id)
	v, err, _ := r.sg.Do(key, func() (interface{}, error) {
		var bus domain.Bus
		if err := r.db.WithContext(ctx).First(&bus, id).Error; err != nil {
			return nil, err
		}
		return &bus, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*domain.Bus), nil
}

func (r *mysqlBusRepository) List(ctx context.Context, origin, dest string, date string) ([]*domain.Bus, error) {
	query := r.db.WithContext(ctx).Model(&domain.Bus{})

	if origin != "" {
		query = query.Where("origin = ?", origin)
	}
	if dest != "" {
		query = query.Where("dest = ?", dest)
	}
	if date != "" {
		query = query.Where("DATE(start_time) = ?", date)
	}

	var buses []*domain.Bus
	if err := query.Order("start_time ASC").Find(&buses).Error; err != nil {
		return nil, err
	}
	return buses, nil
}

func (r *mysqlBusRepository) UpdateSeat(ctx context.Context, busID uint64, delta int) error {
	// delta can be negative for deduction
	return r.db.WithContext(ctx).Model(&domain.Bus{}).
		Where("id = ? AND left_seat + ? >= 0", busID, delta).
		UpdateColumn("left_seat", gorm.Expr("left_seat + ?", delta)).Error
}

func (r *mysqlBusRepository) DecrSeat(ctx context.Context, busID uint64) (bool, error) {
	result := r.db.WithContext(ctx).Model(&domain.Bus{}).
		Where("id = ? AND left_seat > 0", busID).
		UpdateColumn("left_seat", gorm.Expr("left_seat - 1"))
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *mysqlBusRepository) IncrSeat(ctx context.Context, busID uint64) error {
	return r.db.WithContext(ctx).Model(&domain.Bus{}).
		Where("id = ?", busID).
		UpdateColumn("left_seat", gorm.Expr("left_seat + 1")).Error
}
