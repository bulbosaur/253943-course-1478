package main

import (
	"fmt"
	"log"
	"lyceum/config"
	"lyceum/internal/storage"
	v1 "lyceum/internal/transport/gRPC"
	"net"
	"path/filepath"

	pb "lyceum/pkg/api/test"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configDir := "./config"
	envPath := filepath.Join(configDir, ".env")
	yamlPath := filepath.Join(configDir, "config.yaml")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg, err := config.LoadConfig(envPath, yamlPath)
	if err != nil {
		logger.Error("failed to load config", zap.String("envPath", envPath), zap.String("yaml.Path", yamlPath), zap.Error(err))
	}
	
	logger.Info("starting debezium", zap.String("version", "test"), zap.Any("config", cfg))
	logger.Info("logger level", zap.Any("level", logger.Level()))

	orderStorage := storage.NewOrderStorage()
	orderService := v1.NewOrderServiceServer(orderStorage)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("main.StartGrpc: failed to listen", zap.String("addr", addr), zap.Error(err))
	}

	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatalf("main.Grpc: %v", err)
	}

	select {}
}
