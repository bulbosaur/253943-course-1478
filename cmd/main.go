package main

import (
	"log"
	"lyceum/internal/storage"
	v1 "lyceum/internal/transport/gRPC"
	"net"

	pb "lyceum/pkg/api/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = ":50051"

func main() {
	orderStorage := storage.NewOrderStorage()
	orderService := v1.NewOrderServiceServer(orderStorage)

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, orderService)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("main.StartGrpc: %v", err)
	}

	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatalf("main.Grpc: %v", err)
	}

	log.Printf("server started successfully")
}
