package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func New(storagePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err.Error())
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to turn on foreign keys: %s", err.Error())
	}

	return db, nil
}
