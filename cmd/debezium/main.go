package main

import (
	"context"
	"lyceum/config"
	"lyceum/internal/repository"
	"lyceum/internal/storage"
	v1 "lyceum/internal/transport/gRPC"
	srv "lyceum/internal/transport/http"
	lg "lyceum/logger"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	pb "lyceum/pkg/api/test"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ConfigDir = "./config"
	EnvPath   = filepath.Join(ConfigDir, ".env")
	YamlPath  = filepath.Join(ConfigDir, "config.yaml")
)

func main() {
	cfg := config.LoadConfig(EnvPath, YamlPath)

	logger := lg.NewLogger(cfg.Env.LogLevel)
	defer logger.Sync()

	ctx := lg.WithRequestID(context.Background(), "")
	ctx = lg.WithLogger(ctx, logger)

	logger.Info(ctx, "starting gRPC server", zap.String("version", "test"), zap.Any("config", cfg.GRPC))

	orderRepo := repository.StartPostgres(cfg.PostgreSQL)
	logger.Debug(ctx, "successful connection to the database", zap.Any("config", cfg.PostgreSQL))

	redisClient := storage.StartRedisClient(ctx, cfg.Redis)
	logger.Debug(ctx, "successful connection to the redis", zap.Any("config", cfg.Redis))

	orderService := v1.NewOrderServiceServer(orderRepo, redisClient)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(v1.LoggingUnaryInterceptor(logger)),
	)

	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	grpcAddr := net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port))

	lc := &net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", grpcAddr)
	if err != nil {
		logger.Error(ctx, "main.StartGrpc: failed to listen", zap.String("addr", grpcAddr), zap.Error(err))
		return
	}

	go srv.RunRest(ctx, cfg.HTTP)

	err = grpcServer.Serve(l)
	if err != nil {
		logger.Error(
			ctx,
			"main.StartGrpc: failed to serve",
			zap.String("host", cfg.HTTP.Host),
			zap.Int("port", cfg.HTTP.Port),
			zap.Error(err),
		)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Info(ctx, "shutting down gRPC server")
	grpcServer.GracefulStop()
}
