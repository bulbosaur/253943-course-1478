package v1

import (
	"context"
	"fmt"
	"lyceum/internal/models"
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

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *models.CreateOrderRequest) (*models.CreateOrderResponse, error) {
	var resp models.CreateOrderResponse

	orderID := s.storage.CreateOrder(req.Item, req.Quantity)

	resp.ID = orderID
	
	return &resp, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *models.GetOrderRequest) (*models.GetOrderResponse, error) {
	var resp models.GetOrderResponse

	order, err := s.storage.GetOrder(req.ID)
	if err != nil {
		return &models.GetOrderResponse{}, fmt.Errorf("gRPC.GetOrder: %w", err)
	}
	resp.Order = order

	return &resp, nil
}

func (s *OrderServiceServer) UpdateOrder(ctx context.Context, req *models.UpdateOrderRequest) (models.UpdateOrderResponse, error) {
	var resp models.UpdateOrderResponse

	newOrder := s.storage.UpdateOrder(req.ID, req.Item, req.Quantity)
	resp.Order = newOrder

	return resp, nil
}

func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *models.DeleteOrderRequest) (models.DeleteOrderResponse, error) {
	var (
		resp models.DeleteOrderResponse
		err error
	)

	res := s.storage.DeleteOrder(req.ID)
	resp.Success = res

	if !res {
		err = fmt.Errorf("qRPC.DeleteOrder: can't delete an order ID %s", req.ID)
	}

	return resp, err
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *models.ListOrdersRequest) (models.ListOrdersResponse, error) {
	var resp models.ListOrdersResponse

	resp.Orders = s.storage.ListOrders()

	return resp, nil
}