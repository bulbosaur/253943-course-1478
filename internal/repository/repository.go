package repository

import (
	"context"
	"log"
	pb "lyceum/pkg/api/test"
	"lyceum/pkg/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, item string, quantity int32) (string, error)
	GetOrder(ctx context.Context, id string) (*pb.Order, error)
	UpdateOrder(ctx context.Context, id, item string, quantity int32) (*pb.Order, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*pb.Order, error)
}

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{pool: pool}
}

func StartPostgres(cfg db.Config) *PostgresOrderRepository {
	orderPostgres, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("repository.StartPostgres: failed to create db: %d", err)
	}

	orderRepo := NewPostgresOrderRepository(orderPostgres.Pool)
	return orderRepo
}
