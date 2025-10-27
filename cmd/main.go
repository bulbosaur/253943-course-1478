package main

import (
	"context"
	"fmt"
	"log"
	"lyceum/config"
	"lyceum/internal/storage"
	v1 "lyceum/internal/transport/gRPC"
	lg "lyceum/logger"
	"net"
	"path/filepath"

	pb "lyceum/pkg/api/test"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		configDir = "./config"
		envPath = filepath.Join(configDir, ".env")
		yamlPath = filepath.Join(configDir, "config.yaml")
	)

	cfg, err := config.LoadConfig(envPath, yamlPath)
	if err != nil {
		log.Print("failed to load config:", err)
	}

	logger, _ := lg.NewLogger(cfg.Env.LogLevel)
	lg.GlobalLogger = logger
	
	ctx := lg.WithRequestID(context.Background(), "")
	ctx = lg.WithLogger(ctx, logger)
	
	logger.Info(ctx, "starting debezium", zap.String("version", "test"), zap.Any("config", cfg))

	orderStorage := storage.NewOrderStorage()
	orderService := v1.NewOrderServiceServer(orderStorage)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(v1.LoggingUnaryInterceptor),
	)

	pb.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(ctx, "main.StartGrpc: failed to listen", zap.String("addr", addr), zap.Error(err))
		return
	}


	err = grpcServer.Serve(l)
	if err != nil {
		logger.Error(ctx, "main.StartGrpc: failed to serve", zap.String("addr", addr), zap.Error(err))
	}

	select {}
}
