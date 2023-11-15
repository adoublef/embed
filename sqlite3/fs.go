package sqlite3

import (
	"context"
	"io/fs"

	"github.com/maragudk/migrate"
)

type FS struct {
	fsys fs.FS
}

func NewFS(fsys fs.FS) *FS {
	return &FS{fsys: fsys}
}

func (fsys FS) Up(ctx context.Context, db *DB) (err error) {
	return migrate.Up(ctx, db.rwc, fsys.fsys)
}