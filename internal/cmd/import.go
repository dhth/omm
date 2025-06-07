package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
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
	task := types.Task{
		Summary:   taskSummary,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err = pers.InsertTasks(db, []types.Task{task}, true)
	return err
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
	tasks := make([]types.Task, len(taskSummaries))
	for i, summ := range taskSummaries {
		tasks[i] = types.Task{
			Summary:   summ,
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	_, err = pers.InsertTasks(db, tasks, true)
	return err
}
