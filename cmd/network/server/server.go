package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	service "github.com/adoublef/embed/internal/network/http"
	"github.com/adoublef/embed/nats"
	sql "github.com/adoublef/embed/sqlite3"
	"github.com/adoublef/embed/static"
	"github.com/adoublef/embed/template"
	"github.com/go-chi/chi/v5"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	fs := service.T.Funcs(template.DefaultFuncs, static.FuncMap)
	
	mux := chi.NewMux()
	mux.Mount("/", service.New(fs, db, kv))
	mux.Handle("/static/*", http.StripPrefix("/static/", static.Handler))
	
	s := &http.Server{Addr: addr, Handler: mux}
	return &Server{s}, nil
}
