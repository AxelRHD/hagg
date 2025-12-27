package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// OpenSQLite opens a SQLite database with sane defaults.
func OpenSQLite(path string) (*sqlx.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite path is empty")
	}

	// DSN with recommended pragmas
	dsn := fmt.Sprintf(
		"%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000",
		path,
	)

	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
