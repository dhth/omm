package cmd

import (
	"database/sql"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
)

func importTask(db *sql.DB, taskSummary string) error {
	now := time.Now()
	return pers.ImportTaskIntoDB(db, taskSummary, true, now, now)
}

func importTasks(db *sql.DB, taskSummaries []string) error {
	now := time.Now()
	return pers.ImportTasksIntoDB(db, taskSummaries, true, now, now)
}
