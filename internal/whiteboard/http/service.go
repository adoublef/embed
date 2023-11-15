package http

import (
	"fmt"
	"net/http"

	"github.com/adoublef/mvp/nats"
	sql "github.com/adoublef/mvp/sqlite3"
)

type Service struct {
	m  *http.ServeMux
	kv *nats.KV
	db *sql.DB
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

// A New Service will be created
func New(db *sql.DB, kv *nats.KV) *Service {
	s := Service{
		m: http.NewServeMux(),
		db: db,
		kv: kv,
	}
	s.routes()
	return &s
}

func (s *Service) routes() {
	s.m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "whiteboard")
	})
}
