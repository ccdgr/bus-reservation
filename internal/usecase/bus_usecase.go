package usecase

import (
	"context"
	"github.com/ccdgr/bus-reservation/internal/domain"
)

type busUsecase struct {
	busRepo   domain.BusRepository
	redisRepo domain.BusRepository
}

func NewBusUsecase(busRepo domain.BusRepository, redisRepo domain.BusRepository) domain.BusUsecase {
	return &busUsecase{
		busRepo:   busRepo,
		redisRepo: redisRepo,
	}
}

func (u *busUsecase) List(ctx context.Context, origin, dest string, date string) ([]*domain.Bus, error) {
	buses, err := u.busRepo.List(ctx, origin, dest, date)
	if err != nil {
		return nil, err
	}

	// 并发从 Redis 丰富实时库存 (或简单循环)
	for _, b := range buses {
		if stock, err := u.redisRepo.GetStock(ctx, b.ID); err == nil {
			b.LeftSeat = stock
		}
	}

	return buses, nil
}

func (u *busUsecase) GetByID(ctx context.Context, id uint64) (*domain.Bus, error) {
	bus, err := u.busRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 从 Redis 获取最实时的库存
	if stock, err := u.redisRepo.GetStock(ctx, id); err == nil {
		bus.LeftSeat = stock
	}

	return bus, nil
}
