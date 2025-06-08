package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	pers "github.com/dhth/omm/internal/persistence"
)

var (
	errNoTaskAtIndex            = errors.New("no task exists at given index")
	errCouldntMarshalTaskToJSON = errors.New("couldn't marshall task to JSON")
	errCouldntSearchTasks       = errors.New("couldn't search tasks")
)

func listTasks(db *sql.DB, limit, offset uint16, writer io.Writer) error {
	tasks, err := pers.FetchActiveTasks(db, limit, offset)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	summaries := make([]string, len(tasks))
	for i, task := range tasks {
		summaries[i] = task.Summary
	}

	fmt.Fprintln(writer, strings.Join(summaries, "\n"))
	return nil
}

func showTask(db *sql.DB, index uint64, format taskOutputFormat, writer io.Writer) error {
	task, found, err := pers.FetchNthActiveTask(db, index-1)
	if err != nil {
		return err
	}

	if !found {
		return errNoTaskAtIndex
	}

	taskDetails := task.GetDetails()

	switch format {
	case taskOutputPlain:
		fmt.Fprintln(writer, taskDetails.Summary)
		if taskDetails.Context != nil {
			fmt.Fprintf(writer, "\n%s", *taskDetails.Context)
		}
	case taskOutputJSON:
		data, err := json.MarshalIndent(taskDetails, "", "  ")
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntMarshalTaskToJSON, err.Error())
		}
		fmt.Fprintf(writer, "%s\n", data)
	}

	return nil
}

func searchTasks(db *sql.DB, query string, format taskOutputFormat, active bool, limit, offset uint16, writer io.Writer) error {
	tasks, err := pers.FetchTasksThatMatchQuery(db, query, active, limit, offset)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntSearchTasks, err.Error())
	}

	if len(tasks) == 0 {
		return nil
	}

	switch format {
	case taskOutputPlain:
		summaries := make([]string, len(tasks))
		for i, task := range tasks {
			summaries[i] = task.Summary
		}
		fmt.Fprintln(writer, strings.Join(summaries, "\n"))
	case taskOutputJSON:
		data, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntMarshalTaskToJSON, err.Error())
		}
		fmt.Fprintf(writer, "%s\n", data)
	}

	return nil
}
