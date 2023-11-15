package http

import (
	"fmt"
	"net/http"
)

type Service struct {
	m *http.ServeMux
}

func (s*Service) ServeHTTP(w http.ResponseWriter, r*http.Request) {
	s.m.ServeHTTP(w, r)
}

// A New Service will be created
func New() *Service {
	s := Service{
		m: http.NewServeMux(),
	}
	s.routes()
	return &s
}

func (s *Service) routes() {
	s.m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
}
