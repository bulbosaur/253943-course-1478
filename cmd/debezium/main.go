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
	"time"

	pb "lyceum/pkg/api/test"
	"lyceum/pkg/faulttolerance"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		configDir         = "./config"
		envPath           = filepath.Join(configDir, ".env")
		yamlPath          = filepath.Join(configDir, "config.yaml")
		orderRepo         *repository.PostgresOrderRepository
		redisClient       *storage.RedisOrderCache
		defaultMaxRetries = 5
		defaultBaseDelay  = 100 * time.Millisecond
	)
	cfg := config.LoadConfig(envPath, yamlPath)

	logger := lg.NewLogger(cfg.Env.LogLevel)
	defer logger.Sync() //nolint:errcheck // error checking is redundant here

	ctx := lg.WithRequestID(context.Background(), "")
	ctx = lg.WithLogger(ctx, logger)

	logger.Info(ctx, "starting gRPC server", zap.String("version", "test"), zap.Any("config", cfg.GRPC))

	err := faulttolerance.Retry(func() error {
		orderRepo = repository.StartPostgres(cfg.PostgreSQL)
		return orderRepo.Pool.Ping(ctx)
	}, defaultMaxRetries, defaultBaseDelay)
	if err != nil {
		logger.Error(ctx, "failed to connect to PostgreSQL", zap.Any("error", err))
		return
	}
	logger.Debug(ctx, "successful connection to the database", zap.Any("config", cfg.PostgreSQL))

	err = faulttolerance.Retry(func() error {
		redisClient = storage.StartRedisClient(ctx, cfg.Redis)
		return redisClient.Client.Ping(ctx).Err()
	}, defaultMaxRetries, defaultBaseDelay)
	if err != nil {
		logger.Error(ctx, "failed to connect to Redis", zap.Any("error", err))
		return
	}
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
