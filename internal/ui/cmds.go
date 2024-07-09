package ui

import (
	"time"

	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
	_ "modernc.org/sqlite"
)

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}

func updateTaskSequence(db *sql.DB, sequence []uint64) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTaskSequenceInDB(db, sequence)
		return taskSequenceUpdatedMsg{err}
	}
}

func createTask(db *sql.DB, summary string, createdAt, updatedAt time.Time) tea.Cmd {
	return func() tea.Msg {
		id, err := pers.InsertTaskInDB(db, summary, createdAt, updatedAt)
		return taskCreatedMsg{id, summary, createdAt, updatedAt, err}
	}
}

func deleteTask(db *sql.DB, id uint64, index int, active bool) tea.Cmd {
	return func() tea.Msg {
		err := pers.DeleteTaskInDB(db, id)
		return taskDeletedMsg{id, index, active, err}
	}
}

func updateTaskSummary(db *sql.DB, listIndex int, id uint64, summary string) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTaskSummaryInDB(db, id, summary)
		return taskSummaryUpdatedMsg{listIndex, id, summary, err}
	}
}

func changeTaskStatus(db *sql.DB, listIndex int, id uint64, active bool, updatedAt time.Time) tea.Cmd {
	return func() tea.Msg {
		err := pers.ChangeTaskStatusInDB(db, id, active, updatedAt)
		return taskStatusChangedMsg{listIndex, id, active, updatedAt, err}
	}
}

func fetchTasks(db *sql.DB, active bool, limit int) tea.Cmd {
	return func() tea.Msg {
		var tasks []types.Task
		var err error
		switch active {
		case true:
			tasks, err = pers.FetchActiveTasksFromDB(db, limit)
		case false:
			tasks, err = pers.FetchInActiveTasksFromDB(db, limit)
		}
		return tasksFetched{tasks, active, err}
	}
}
