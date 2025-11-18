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
	Client *redis.Client
}

func NewRedisOrderCache(client *redis.Client) *RedisOrderCache {
	return &RedisOrderCache{Client: client}
}

func (s *RedisOrderCache) GetOrder(ctx context.Context, id string) (*pb.Order, error) {
	val, err := s.Client.Get(ctx, "order:"+id).Bytes()
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

func (s *RedisOrderCache) SetOrder(ctx context.Context, id string, order *pb.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, "order:"+id, data, ttl).Err()
}

func (s *RedisOrderCache) DeleteOrder(ctx context.Context, id string) error {
	return s.Client.Del(ctx, "order:"+id).Err()
}
