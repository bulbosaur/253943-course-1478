package srv

import (
	"context"
	"net/http"
	"strconv"
)

type Server struct {
	srv *http.Server
}

func NewServer(port int) *Server {
	srv := http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: nil,
	}

	return &Server{srv: &srv}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}