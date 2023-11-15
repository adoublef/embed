package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	wbHTTP "github.com/adoublef/mvp/internal/whiteboard/http"
	"github.com/adoublef/mvp/nats"
)

var (
	timeout = 5 * time.Second
)

type Server struct {
	s *http.Server
}

// ListenAndServe listens on the TCP network address and then handles requests on incoming connections.
func (s *Server) ListenAndServe() error {
	return s.s.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.s.Shutdown(ctx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// A New Server defines parameters for running an HTTP server.
func NewServer(nc *nats.Conn) (*Server, error) {
	m := wbHTTP.New()
	s := &http.Server{
		Addr:    ":8080",
		Handler: m,
	}
	return &Server{s}, nil
}
