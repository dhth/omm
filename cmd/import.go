package cmd

import (
	"database/sql"
	"fmt"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
)

var errImportWillExceedTaskLimit = fmt.Errorf("import will exceed maximum number of tasks allowed, which is %d. Archive/Delete tasks that are not active using ctrl+d/ctrl+x", pers.TaskNumLimit)

func importTask(db *sql.DB, taskSummary string) error {
	numTasks, err := pers.FetchNumActiveTasksFromDB(db)
	if err != nil {
		return err
	}
	if numTasks+1 > pers.TaskNumLimit {
		return errImportWillExceedTaskLimit
	}

	now := time.Now()
	return pers.ImportTaskIntoDB(db, taskSummary, true, now, now)
}

func importTasks(db *sql.DB, taskSummaries []string) error {
	numTasks, err := pers.FetchNumActiveTasksFromDB(db)
	if err != nil {
		return err
	}
	if numTasks+len(taskSummaries) > pers.TaskNumLimit {
		return errImportWillExceedTaskLimit
	}

	now := time.Now()
	return pers.ImportTaskSummariesIntoDB(db, taskSummaries, true, now, now)
}
