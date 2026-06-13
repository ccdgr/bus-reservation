package repository

import (
	"context"
	"fmt"
	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisBusRepository struct {
	client *redis.Client
}

func NewRedisBusRepository(client *redis.Client) *RedisBusRepository {
	return &RedisBusRepository{client: client}
}

const decrSeatLua = `
local key = KEYS[1]
local stock = tonumber(redis.call('get', key))
if stock and stock > 0 then
    redis.call('decr', key)
    return 1
else
    return 0
end
`

func (r *RedisBusRepository) DecrSeat(ctx context.Context, busID uint64) (bool, error) {
	key := fmt.Sprintf("bus_stock:%d", busID)
	result, err := r.client.Eval(ctx, decrSeatLua, []string{key}).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (r *RedisBusRepository) IncrSeat(ctx context.Context, busID uint64) error {
	key := fmt.Sprintf("bus_stock:%d", busID)
	return r.client.Incr(ctx, key).Err()
}

func (r *RedisBusRepository) GetStock(ctx context.Context, busID uint64) (int, error) {
	key := fmt.Sprintf("bus_stock:%d", busID)
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisBusRepository) GetByID(ctx context.Context, id uint64) (*domain.Bus, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *RedisBusRepository) List(ctx context.Context, origin, dest string, date string) ([]*domain.Bus, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *RedisBusRepository) UpdateSeat(ctx context.Context, busID uint64, delta int) error {
	return fmt.Errorf("not implemented")
}

func (r *RedisBusRepository) SetStock(ctx context.Context, busID uint64, stock int) error {
	key := fmt.Sprintf("bus_stock:%d", busID)
	return r.client.Set(ctx, key, stock, 0).Err()
}
