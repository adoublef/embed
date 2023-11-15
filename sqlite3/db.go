package sqlite3

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	rwc *sql.DB
}

// Close closes the database and prevents new queries from starting. C
func (d *DB) Close() error {
	return d.rwc.Close()
}

// Open opens a database connection for the given sqlite file
func Open(dsn string) (*DB, error) {
	var (
		args = strings.Join([]string{"_journal=wal", "_timeout=5000", "_synchronous=normal", "_fk=true"}, "&")
	)
	rwc, err := sql.Open("sqlite3", dsn+"?"+args)
	if err != nil {
		return nil, err
	}
	return &DB{rwc}, nil
}
