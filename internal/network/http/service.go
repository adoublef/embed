package http

import (
	"embed"
	"net/http"

	"github.com/adoublef/embed/nats"
	sql "github.com/adoublef/embed/sqlite3"
	"github.com/adoublef/embed/template"
)

var (
	pageIndex = "index.html"
	pagePost  = "post.html"
)

//go:embed all:*.html
var embedFS embed.FS
var T = template.NewFS(embedFS)

type Service struct {
	m  *http.ServeMux
	kv *nats.KV
	db *sql.DB
	fs *template.FS
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

// A New Service will be created
func New(fs *template.FS, db *sql.DB, kv *nats.KV) *Service {
	s := Service{
		m:  http.NewServeMux(),
		fs: fs,
		db: db,
		kv: kv,
	}
	s.routes()
	return &s
}

func (s *Service) routes() {
	s.m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.fs.ParseFiles(pageIndex)
		if err != nil {
			http.Error(w, "Failed to parse", http.StatusUnprocessableEntity)
			return 
		}

		t.Execute(w, nil)
	})

	s.m.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.fs.ParseFiles(pageIndex, pagePost)
		if err != nil {
			http.Error(w, "Failed to parse", http.StatusUnprocessableEntity)
			return 
		}
		t.Execute(w, nil)
	})
}
