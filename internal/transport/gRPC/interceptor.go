package v1

import (
	"context"
	"lyceum/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	var reqID string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	if listID := md["x-request-id"]; len(listID) > 0 {
		reqID = listID[0]
	} else {
		reqID = uuid.NewString()
	}

	ctx = logger.WithRequestID(ctx, reqID)
	ctx = logger.WithLogger(ctx, logger.GlobalLogger)

	return handler(ctx, req)
}