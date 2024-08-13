package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
)

var errWillExceedCapacity = errors.New("import will exceed capacity")

func importTask(db *sql.DB, taskSummary string) error {
	numTasks, err := pers.FetchNumActiveTasksShown(db)
	if err != nil {
		return err
	}
	if numTasks+1 > pers.TaskNumLimit {
		return fmt.Errorf("%w (current task count: %d)", errWillExceedCapacity, numTasks)
	}

	now := time.Now()
	return pers.ImportTask(db, taskSummary, true, now, now)
}

func importTasks(db *sql.DB, taskSummaries []string) error {
	numTasks, err := pers.FetchNumActiveTasksShown(db)
	if err != nil {
		return err
	}
	if numTasks+len(taskSummaries) > pers.TaskNumLimit {
		return fmt.Errorf("%w (current task count: %d)", errWillExceedCapacity, numTasks)
	}

	now := time.Now()
	return pers.ImportTaskSummaries(db, taskSummaries, true, now, now)
}
