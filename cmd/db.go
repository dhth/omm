package cmd

import (
	"database/sql"
	"time"
)

func getDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbpath)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, err
}

func initDB(db *sql.DB) error {
	// these init queries cannot be changed once omm is released; only further
	// migrations can be added, which are run when omm sees a difference between
	// the values in the db_versions table and latestDBVersion
	_, err := db.Exec(`
CREAT TABLE IF NOT EXISTS db_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    summary TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE task_sequence (
    id INTEGER PRIMARY KEY,
    sequence JSON NOT NULL
);

INSERT INTO task_sequence (id, sequence) VALUES (1, '[]');

INSERT INTO db_versions (version, created_at)
VALUES (1, ?);
`, time.Now().UTC())

	return err
}
