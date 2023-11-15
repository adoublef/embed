package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	service "github.com/adoublef/mvp/internal/network/http"
	"github.com/adoublef/mvp/nats"
	sql "github.com/adoublef/mvp/sqlite3"
)

type Server struct {
	s *http.Server
}

// ListenAndServe listens on the TCP network address and then handles requests on incoming connections.
func (s *Server) ListenAndServe() error {
	if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http listen and serve: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := s.s.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http shutdown: %w", err)
	}
	return nil
}

// A New Server defines parameters for running an HTTP server.
func NewServer(addr string, nc *nats.Conn, db *sql.DB) (*Server, error) {
	jsc, err := nats.JetStream(nc)
	if err != nil {
		return nil, err
	}
	kv, err := nats.UpsertKV(jsc, &nats.KVConfig{
		Bucket: "temp",
	})
	if err != nil {
		return nil, err
	}
	m := service.New(db, kv)
	s := &http.Server{Addr: addr, Handler: m}
	return &Server{s}, nil
}
