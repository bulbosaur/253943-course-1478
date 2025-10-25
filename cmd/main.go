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

	cfg, err := config.LoadConfig(envPath, yamlPath)
	if err != nil { return }
	

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("starting debezium", zap.String("version", "test"), zap.Any("config", cfg))

	orderStorage := storage.NewOrderStorage()
	orderService := v1.NewOrderServiceServer(orderStorage)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("main.StartGrpc: %v", err)
	}

	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatalf("main.Grpc: %v", err)
	}

	logger.Info("server started successfully")

	select {}
}
