package storage

import (
	"context"
	"encoding/json"
	"time"

	pb "lyceum/pkg/api/test"

	"github.com/redis/go-redis/v9"
)

type OrderCache interface {
	GetOrder(ctx context.Context, id string) (*pb.Order, error)
	SetOrder(ctx context.Context, id string, order *pb.Order, ttl time.Duration) error
	DeleteOrder(ctx context.Context, id string) error
}

type RedisOrderCache struct {
	client *redis.Client
}

func NewRedisOrderCache(client *redis.Client) *RedisOrderCache {
	return &RedisOrderCache{client: client}
}

func (r *RedisOrderCache) GetOrder(ctx context.Context, id string) (*pb.Order, error) {
	val, err := r.client.Get(ctx, "order:"+id).Bytes()
	if err != nil {
		return nil, err
	}

	var order pb.Order
	err = json.Unmarshal(val, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RedisOrderCache) SetOrder(ctx context.Context, id string, order *pb.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "order:"+id, data, ttl).Err()
}

func (r *RedisOrderCache) DeleteOrder(ctx context.Context, id string) error {
	return r.client.Del(ctx, "order:"+id).Err()
}
