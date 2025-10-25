package v1

// import (
// 	"context"

// 	lg "lyceum/logger"

// 	"google.golang.org/grpc"
// )

// func RegistLoggerInterceptor(logger lg.Logger) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
// 		ctx = lg.WithLogger(ctx, logger)
// 		return handler(ctx, req)
// 	}
// }

