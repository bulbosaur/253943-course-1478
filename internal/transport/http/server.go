package srv

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
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
