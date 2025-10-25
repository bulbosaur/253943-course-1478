package v1

import (
	"context"

	"lyceum/logger"

	"google.golang.org/grpc"
)

func RegistLoggerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = logger.WithLogger(ctx, logger)
		return handler(ctx, req)
	}
}
