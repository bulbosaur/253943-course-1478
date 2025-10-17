package v1

import (
	"lyceum/internal/storage"
	pb "lyceum/pkg/api/test"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	storage *storage.OrderStorage
}

func NewOrderServiceServer(storage *storage.OrderStorage) *OrderServiceServer {
	return &OrderServiceServer{
		storage: storage,
	}
}