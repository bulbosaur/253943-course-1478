package v1

import (
	"lyceum/internal/repository"
	pb "lyceum/pkg/api/test"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer

	repository repository.OrderRepository
}

func NewOrderServiceServer(repo repository.OrderRepository) *OrderServiceServer {
	return &OrderServiceServer{
		repository: repo,
	}
}
