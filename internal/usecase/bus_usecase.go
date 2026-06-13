package usecase

import (
	"context"
	"github.com/ccdgr/bus-reservation/internal/domain"
)

type busUsecase struct {
	busRepo domain.BusRepository
}

func NewBusUsecase(busRepo domain.BusRepository) domain.BusUsecase {
	return &busUsecase{busRepo: busRepo}
}

func (u *busUsecase) List(ctx context.Context, origin, dest string, date string) ([]*domain.Bus, error) {
	return u.busRepo.List(ctx, origin, dest, date)
}

func (u *busUsecase) GetByID(ctx context.Context, id uint64) (*domain.Bus, error) {
	return u.busRepo.GetByID(ctx, id)
}
