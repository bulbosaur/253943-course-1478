package v1

import (
	"context"
	"errors"
	"fmt"
	"lyceum/logger"
	pb "lyceum/pkg/api/test"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderServiceServer) CreateOrder(
	ctx context.Context,
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
	var resp pb.CreateOrderResponse

	if req.GetItem() == "" {
		return &resp, status.Error(codes.InvalidArgument, "gRPC.CreateOrder: item is required")
	}

	if req.GetQuantity() <= 0 {
		return &resp, status.Error(codes.InvalidArgument, "gRPC.CreateOrder: quantity must be positive")
	}

	orderID := s.storage.CreateOrder(req.GetItem(), req.GetQuantity())

	l := logger.FromContext(ctx)
	l.Debug(ctx, "new order was created", zap.String("orderID", orderID))

	resp.Id = orderID

	return &resp, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var resp pb.GetOrderResponse

	order, err := s.storage.GetOrder(req.GetId())

	l := logger.FromContext(ctx)
	
	if err != nil {
		l.Error(ctx, "gRPC.GetOrder", zap.Any("error", err))
		return &pb.GetOrderResponse{}, fmt.Errorf("gRPC.GetOrder: %w", err)
	}
	l.Debug(ctx, "order was got", zap.Any("order", order))

	resp.Order = order

	return &resp, nil
}

func (s *OrderServiceServer) UpdateOrder(
	ctx context.Context,
	req *pb.UpdateOrderRequest,
) (*pb.UpdateOrderResponse, error) {
	var resp pb.UpdateOrderResponse

	if req.GetId() == "" {
		return &resp, fmt.Errorf("gRPC.UpdateOrder: %w", errors.New("orderID is empty"))
	}

	newOrder := s.storage.UpdateOrder(req.GetId(), req.GetItem(), req.GetQuantity())
	resp.Order = newOrder

	l := logger.FromContext(ctx)
	l.Debug(ctx, "order was updated", zap.Any("newOrder", newOrder))

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

	res := s.storage.DeleteOrder(req.GetId())
	resp.Success = res

	if !res {
		err = fmt.Errorf("qRPC.DeleteOrder: can't delete an order ID %s", req.GetId())
	}

	l := logger.FromContext(ctx)
	l.Debug(ctx, "order was updated", zap.String("orderID", req.GetId()))

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
