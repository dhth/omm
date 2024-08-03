package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	latestDBVersion = 2 // only upgrade this after adding a migration in getMigrations
)

var (
	errDBDowngraded          = errors.New("database downgraded")
	errDBMigrationFailed     = errors.New("database migration failed")
	errCouldntFetchDBVersion = errors.New("couldn't fetch version")
)

type dbVersionInfo struct {
	id        int
	version   int
	createdAt time.Time
}

func getMigrations() map[int]string {
	migrations := make(map[int]string)
	// these migrations should not be modified once released.
	// that is, migrations is an append-only map.

	migrations[2] = `
ALTER TABLE task
ADD COLUMN context TEXT;
`

	return migrations
}

func fetchLatestDBVersion(db *sql.DB) (dbVersionInfo, error) {
	row := db.QueryRow(`
SELECT id, version, created_at
FROM db_versions
ORDER BY created_at DESC
LIMIT 1;
`)

	var dbVersion dbVersionInfo
	err := row.Scan(
		&dbVersion.id,
		&dbVersion.version,
		&dbVersion.createdAt,
	)

	return dbVersion, err
}

func upgradeDBIfNeeded(db *sql.DB) error {
	latestVersionInDB, versionErr := fetchLatestDBVersion(db)
	if versionErr != nil {
		return fmt.Errorf("%w: %w", errCouldntFetchDBVersion, versionErr)
	}

	if latestVersionInDB.version > latestDBVersion {
		return errDBDowngraded
	}

	if latestVersionInDB.version < latestDBVersion {
		err := upgradeDB(db, latestVersionInDB.version)
		if err != nil {
			return err
		}
	}

	return nil
}

func upgradeDB(db *sql.DB, currentVersion int) error {
	migrations := getMigrations()
	for i := currentVersion + 1; i <= latestDBVersion; i++ {
		migrateQuery := migrations[i]
		migrateErr := runMigration(db, migrateQuery, i)
		if migrateErr != nil {
			return fmt.Errorf("%w (version %d): %w", errDBMigrationFailed, i, migrateErr)
		}
	}
	return nil
}

func runMigration(db *sql.DB, migrateQuery string, version int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(migrateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
INSERT INTO db_versions (version, created_at)
VALUES (?, ?);
`)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(version, time.Now().UTC())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
