package v1

import (
	"context"
	"lyceum/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = logger.WithRequestID(ctx, uuid.NewString())
		ctx = logger.WithLogger(ctx, log)
		return handler(ctx, req)
	}
}