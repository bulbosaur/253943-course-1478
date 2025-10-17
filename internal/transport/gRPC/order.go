package v1

import (
	storage "lyceum/internal/storage"
	pb "lyceum/pkg/api/test"
)
type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	storage *storage.OrderStorage
}

func NewOrderServiceServer() *OrderServiceServer {
	return &OrderServiceServer{
		pb.UnimplementedOrderServiceServer{},
		storage.NewOrderStorage(),
	}
}

func (s *OrderServiceServer) CreateOrder()