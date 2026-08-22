package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const DefaultPath = "./data/synclink.db"

func PathFromEnv() string {
	if p := os.Getenv("SYNCLINK_DB"); p != "" {
		return p
	}
	return DefaultPath
}

func Open(path string) (*sql.DB, error) {
	if path == "" {
		path = PathFromEnv()
	}
	dsn := path
	if path == ":memory:" {
		dsn = "file:synclink?mode=memory&cache=shared"
	} else {
		dir := filepath.Dir(path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, err
			}
		}
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := Migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
