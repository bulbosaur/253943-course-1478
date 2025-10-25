package v1

import (
	"context"
	"fmt"
	pb "lyceum/pkg/api/test"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderServiceServer) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	var resp pb.CreateOrderResponse

	if req.Item == "" {
		return &resp, status.Error(codes.InvalidArgument, "gRPC.CreateOrder: item is required")
	}

	if req.Quantity <= 0 {
		return &resp, status.Error(codes.InvalidArgument, "gRPC.CreateOrder: quantity must be positive")
	}

	orderID := s.storage.CreateOrder(req.Item, req.Quantity)

	resp.Id = orderID

	return &resp, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var resp pb.GetOrderResponse

	order, err := s.storage.GetOrder(req.Id)
	if err != nil {
		return &pb.GetOrderResponse{}, fmt.Errorf("gRPC.GetOrder: %w", err)
	}
	resp.Order = order

	return &resp, nil
}

func (s *OrderServiceServer) UpdateOrder(
	ctx context.Context,
	req *pb.UpdateOrderRequest,
) (*pb.UpdateOrderResponse, error) {
	var resp pb.UpdateOrderResponse

	newOrder := s.storage.UpdateOrder(req.Id, req.Item, req.Quantity)
	resp.Order = newOrder

	return &resp, nil
}

func (s *OrderServiceServer) DeleteOrder(
	ctx context.Context,
	req *pb.DeleteOrderRequest,
) (*pb.DeleteOrderResponse, error) {
	var (
		resp pb.DeleteOrderResponse
		err  error
	)

	res := s.storage.DeleteOrder(req.Id)
	resp.Success = res

	if !res {
		err = fmt.Errorf("qRPC.DeleteOrder: can't delete an order ID %s", req.Id)
	}

	return &resp, err
}

func (s *OrderServiceServer) ListOrders(
	ctx context.Context,
	req *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
	var resp pb.ListOrdersResponse

	resp.Orders = s.storage.ListOrders()

	return &resp, nil
}
