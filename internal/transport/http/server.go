package srv

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const defaultReadHeaderTimeout time.Duration = 5

type Server struct {
	srv *http.Server
}

func NewServer(port int) *Server {
	srv := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	return &Server{srv: &srv}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
