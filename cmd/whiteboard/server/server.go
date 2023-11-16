package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	service "github.com/adoublef/embed/internal/whiteboard/http"
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
	return s.s.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := s.s.Shutdown(ctx)
	if err != nil || !errors.Is(err, http.ErrServerClosed) {
		return err
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
	t, err := service.T.Funcs(template.DefaultFuncs, static.FuncMap).Parse();
	if  err != nil {
		return nil, err
	}
	h := chi.NewMux()
	h.Mount("/", service.New(t, db, kv))
	h.Handle("/assets/*", http.StripPrefix("/assets/", static.Handler))
	s := &http.Server{Addr: addr, Handler: h}
	return &Server{s}, nil
}


/* 
mux.Handle("/", http.RedirectHandler("/feed", http.StatusFound))
	mux.Mount("/feed", service.New(t, db, nc))
*/