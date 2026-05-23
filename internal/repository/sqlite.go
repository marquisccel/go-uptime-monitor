package repository

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLite(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return db, runMigrations(db)
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS targets (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		name       TEXT NOT NULL,
		url        TEXT NOT NULL UNIQUE,
		interval   INTEGER NOT NULL DEFAULT 60,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS check_results (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		target_id   INTEGER NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
		status_code INTEGER,
		latency_ms  INTEGER,
		is_up       BOOLEAN NOT NULL,
		error       TEXT,
		checked_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_cr_target ON check_results(target_id);
	CREATE INDEX IF NOT EXISTS idx_cr_time   ON check_results(checked_at);
	`)
	return err
}
