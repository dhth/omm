package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"

	pers "github.com/dhth/omm/internal/persistence"
)

var errNoTaskAtIndex = errors.New("no task exists at given index")

func printTasks(db *sql.DB, limit uint16, writer io.Writer) error {
	tasks, err := pers.FetchActiveTasks(db, int(limit))
	if err != nil {
		return err
	}
	summaries := make([]string, len(tasks))
	for i, task := range tasks {
		summaries[i] = task.Summary
	}

	fmt.Fprintln(writer, strings.Join(summaries, "\n"))
	return nil
}

func showTask(db *sql.DB, index uint64, writer io.Writer) error {
	task, found, err := pers.FetchNthActiveTask(db, index)
	if err != nil {
		return err
	}

	if !found {
		return errNoTaskAtIndex
	}

	fmt.Fprintln(writer, task.Summary)
	if task.Context != nil {
		fmt.Fprintf(writer, "\n%s", *task.Context)
	}
	return nil
}
