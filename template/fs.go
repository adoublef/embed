package template

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
)

type FS struct {
	fsys    fs.FS
	pattern []string
}

// Parse reads from the file system fs 
func (fsys *FS) Parse() (Template, error) {
	return template.New("").ParseFS(fsys.fsys, fsys.pattern...)
}

// NewFS
func NewFS(fsys fs.FS, pattern ...string) *FS {
	return &FS{fsys: fsys, pattern: pattern}
}

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

func ExecuteHTTP(w http.ResponseWriter, t Template, name string, data any) {
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusUnprocessableEntity)
	}
}

var DefaultFuncs = template.FuncMap{
	"env": func(s string)string {
		return os.Getenv(s)
	},
} 