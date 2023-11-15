package http

import (
	"embed"
	"net/http"

	"github.com/adoublef/embed/nats"
	sql "github.com/adoublef/embed/sqlite3"
	t "github.com/adoublef/embed/template"
)

//go:embed all:*.html
var embedFS embed.FS
var T = t.NewFS(embedFS, "*.html")

type Service struct {
	m  *http.ServeMux
	kv *nats.KV
	db *sql.DB
	t  t.Template
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

// A New Service will be created
func New(t t.Template, db *sql.DB, kv *nats.KV) *Service {
	s := Service{
		m:  http.NewServeMux(),
		t:  t,
		db: db,
		kv: kv,
	}
	s.routes()
	return &s
}

func (s *Service) routes() {
	s.m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteHTTP(w, s.t, "index.html", nil)
	})
}
