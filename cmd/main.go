package main

import (
	"log"
	"lyceum/internal/storage"
	v1 "lyceum/internal/transport/gRPC"
	"net"

	pb "lyceum/pkg/api/test"

	"google.golang.org/grpc"
)

func main() {
	orderStorage := storage.NewOrderStorage()
	orderService := v1.NewOrderServiceServer(orderStorage)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, orderService)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("main.StartGrpc: %v", err)
	}
	
	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatalf("main.Grpc: %v", err)
	}
}