package template

import (
	"html/template"
	"io"
	"io/fs"
	"maps"
	"net/http"
	"os"
)

type FS struct {
	fsys    fs.FS
	pattern []string
	funcs   template.FuncMap
}

// Parse reads from the file system fs
func (fsys *FS) Parse() (Template, error) {
	return template.New("").Funcs(fsys.funcs).ParseFS(fsys.fsys, fsys.pattern...)
}

// Funcs adds the elements of the argument map to the template's function map. 
func (fsys *FS) Funcs(funcs ...map[string]any) *FS {
	for _, f := range funcs {
		maps.Copy(fsys.funcs, f)
	}
	return fsys
}

// NewFS
func NewFS(fsys fs.FS, pattern ...string) *FS {
	return &FS{fsys: fsys, pattern: pattern, funcs: template.FuncMap{}}
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
	"env": func(s string) string {
		return os.Getenv(s)
	},
}
