package main

import (
	"context"
	"fmt"
	"log"
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

	pb "lyceum/pkg/api/test"
	"lyceum/pkg/db"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
		configDir = "./config"
		envPath   = filepath.Join(configDir, ".env")
		yamlPath  = filepath.Join(configDir, "config.yaml")
	)

func main() {
	cfg, err := config.LoadConfig(envPath, yamlPath)
	if err != nil {
		log.Print("failed to load config:", err)
	}

	logger, err := lg.NewLogger(cfg.Env.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	ctx := lg.WithRequestID(context.Background(), "")
	ctx = lg.WithLogger(ctx, logger)

	logger.Info(ctx, "starting gRPC server", zap.String("version", "test"), zap.Any("config", cfg.GRPC))

	orderPostgres, err := db.NewPostgres(cfg.PostgreSQL)
	if err != nil {
		logger.Error(
			ctx,
			"main.NewPostgres: failed to create db",
			zap.String("addr", fmt.Sprintf("%s:%d", cfg.PostgreSQL.Host, cfg.GRPC.Port)),
			zap.Error(err),
		)
		os.Exit(1)
	}
	logger.Debug(ctx, "successful connection to the database", zap.Any("config", cfg.PostgreSQL))

	orderRepo := repository.NewPostgresOrderRepository(orderPostgres.Pool)

	orderRedis, err := storage.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Error(
			ctx,
			"main.NewRedisClient: failed to create redis client",
			zap.String("addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)),
			zap.Error(err),
		)
		os.Exit(1)
	}
	logger.Debug(ctx, "successful connection to the redis", zap.Any("config", cfg.Redis))

	orderCache := storage.NewRedisOrderCache(orderRedis)

	orderService := v1.NewOrderServiceServer(orderRepo, orderCache)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(v1.LoggingUnaryInterceptor(logger)),
	)

	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	lc := &net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", grpcAddr)
	if err != nil {
		logger.Error(ctx, "main.StartGrpc: failed to listen", zap.String("addr", grpcAddr), zap.Error(err))
		return
	}

	httpAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	httpTimeout := cfg.HTTP.Timeout
	go srv.RunRest(ctx, httpAddr, httpTimeout)

	err = grpcServer.Serve(l)
	if err != nil {
		logger.Error(ctx, "main.StartGrpc: failed to serve", zap.String("addr", httpAddr), zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Info(ctx, "shutting down gRPC server")
	grpcServer.GracefulStop()
}
