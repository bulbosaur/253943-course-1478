package v1

import (
	"context"
	"errors"
	"fmt"
	"lyceum/logger"
	pb "lyceum/pkg/api/test"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultCacheTimeLife = 10 * time.Minute

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

	l := logger.FromContext(ctx)

	orderID, err := s.repository.CreateOrder(ctx, req.GetItem(), req.GetQuantity())
	if err != nil {
		l.Error(ctx, "gRPC.CreateOrder", zap.Any("error", err))
		return nil, fmt.Errorf("gRPC.CreateOrder: %w", err)
	}

	l.Debug(ctx, "new order was created", zap.String("orderID", orderID))

	resp.Id = orderID

	return &resp, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var resp pb.GetOrderResponse
	l := logger.FromContext(ctx)
	id := req.GetId()

	order, err := s.cache.GetOrder(ctx, id)
	if err != nil {
		order, err = s.repository.GetOrder(ctx, id)
		if err != nil {
			l.Error(ctx, "gRPC.GetOrder", zap.Any("error", err))
			return &pb.GetOrderResponse{}, fmt.Errorf("gRPC.GetOrder: %w", err)
		}
	}
	l.Debug(ctx, "order was got", zap.Any("order", order))

	_ = s.cache.SetOrder(ctx, id, order, DefaultCacheTimeLife)

	resp.Order = order

	return &resp, nil
}

func (s *OrderServiceServer) UpdateOrder(
	ctx context.Context,
	req *pb.UpdateOrderRequest,
) (*pb.UpdateOrderResponse, error) {
	var resp pb.UpdateOrderResponse

	l := logger.FromContext(ctx)
	id := req.GetId()
	item := req.GetItem()
	quantity := req.GetQuantity()

	if id == "" {
		return &resp, fmt.Errorf("gRPC.UpdateOrder: %w", errors.New("orderID is empty"))
	}

	newOrder, err := s.repository.UpdateOrder(ctx, id, item, quantity)
	if err != nil {
		l.Error(ctx, "gRPC.UpdateOrder", zap.Any("error", err))
		return nil, fmt.Errorf("gRPC.UpdateOrder: %w", err)
	}

	resp.Order = newOrder

	_ = s.cache.DeleteOrder(ctx, id)

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

	id := req.GetId()
	l := logger.FromContext(ctx)

	res, err := s.repository.DeleteOrder(ctx, id)
	resp.Success = res

	if !res {
		l.Error(ctx, "gRPC.DeleteOrder", zap.Any("error", err))
		err = fmt.Errorf("gRPC.DeleteOrder: can't delete an order ID %s", id)
		return nil, err
	}

	_ = s.cache.DeleteOrder(ctx, id)

	l.Debug(ctx, "order was deletes", zap.String("orderID", id))

	return &resp, err
}

func (s *OrderServiceServer) ListOrders(
	ctx context.Context,
	_ *pb.ListOrdersRequest,
) (*pb.ListOrdersResponse, error) {
	var (
		resp pb.ListOrdersResponse
		err  error
	)

	l := logger.FromContext(ctx)

	resp.Orders, err = s.repository.ListOrders(ctx)
	if err != nil {
		l.Error(ctx, "gRPC.ListOrder", zap.Any("error", err))
		return nil, fmt.Errorf("gRPC.ListOder: %w", err)
	}

	return &resp, nil
}
