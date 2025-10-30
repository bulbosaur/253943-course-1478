package srv

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"lyceum/logger"
	pb "lyceum/pkg/api/test"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultReadHeaderTimeout time.Duration = 5

type Server struct {
	srv *http.Server
}

func NewServer(port int) *Server {
	readHeaderTimeout := viper.GetDuration("HTTP_TIMEOUT") // Проверить, корректно ли
	println(readHeaderTimeout)
	
	srv := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Server{srv: &srv}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func RunRest(ctx context.Context, addr string, timeout time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
	}

	l := logger.FromContext(ctx)
	l.Info(ctx, "starting HTTP server", zap.String("version", "test"), zap.Any("addr", addr), zap.Any("timeout", timeout))
	
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Error(ctx, "http.RunRest: HTTP server failed", zap.Error(err))
		panic(err)
	}
}