package srv

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "lyceum/pkg/api/test"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
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

func RunRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		panic(err)
	}
	log.Printf("server listening at 8081")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}