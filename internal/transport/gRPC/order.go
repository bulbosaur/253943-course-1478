package v1

import (
	storage "gitlab.crja72.ru/golang/2025/spring/course/students/253943-Sofiytula71-gmail.com-course-1478/-/tree/main"
	pb "gitlab.crja72.ru/golang/2025/spring/course/students/253943-Sofiytula71-gmail.com-course-1478/-/tree/main/pkg/api/test"
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