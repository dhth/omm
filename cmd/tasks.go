package cmd

import (
	"database/sql"
	"fmt"
	"io"

	pers "github.com/dhth/omm/internal/persistence"
)

func printTasks(db *sql.DB, limit uint8, writer io.Writer) error {
	tasks, err := pers.FetchActiveTasks(db, int(limit))
	if err != nil {
		return err
	}
	for _, task := range tasks {
		fmt.Fprintf(writer, "%s\n", task.Summary)
	}
	return nil
}
