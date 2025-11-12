package v1

import (
	"lyceum/internal/repository"
	"lyceum/internal/storage"
	pb "lyceum/pkg/api/test"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer

	repository repository.OrderRepository
	cache      storage.OrderCache
}

func NewOrderServiceServer(repo repository.OrderRepository, cache storage.OrderCache) *OrderServiceServer {
	return &OrderServiceServer{
		repository: repo,
		cache:      cache,
	}
}
