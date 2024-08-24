package ui

import (
	"database/sql"
	"os/exec"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
	_ "modernc.org/sqlite" // sqlite driver
)

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}

func updateTaskSequence(db *sql.DB, sequence []uint64) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTaskSequence(db, sequence)
		return taskSequenceUpdatedMsg{err}
	}
}

func createTask(db *sql.DB, summary string, createdAt, updatedAt time.Time) tea.Cmd {
	return func() tea.Msg {
		id, err := pers.InsertTask(db, summary, createdAt, updatedAt)
		return taskCreatedMsg{id, summary, createdAt, updatedAt, err}
	}
}

func deleteTask(db *sql.DB, id uint64, index int, active bool) tea.Cmd {
	return func() tea.Msg {
		err := pers.DeleteTask(db, id)
		return taskDeletedMsg{id, index, active, err}
	}
}

func updateTaskSummary(db *sql.DB, listIndex int, id uint64, summary string) tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		err := pers.UpdateTaskSummary(db, id, summary, now)
		return taskSummaryUpdatedMsg{listIndex, id, summary, now, err}
	}
}

func updateTaskContext(db *sql.DB, listIndex int, id uint64, context string, list taskListType) tea.Cmd {
	return func() tea.Msg {
		var err error
		now := time.Now()
		if context == "" {
			err = pers.UnsetTaskContext(db, id, now)
		} else {
			err = pers.UpdateTaskContext(db, id, context, now)
		}
		return taskContextUpdatedMsg{listIndex, list, id, context, now, err}
	}
}

func changeTaskStatus(db *sql.DB, listIndex int, id uint64, active bool, updatedAt time.Time) tea.Cmd {
	return func() tea.Msg {
		err := pers.ChangeTaskStatus(db, id, active, updatedAt)
		return taskStatusChangedMsg{listIndex, id, active, updatedAt, err}
	}
}

func fetchTasks(db *sql.DB, active bool, limit int) tea.Cmd {
	return func() tea.Msg {
		var tasks []types.Task
		var err error
		switch active {
		case true:
			tasks, err = pers.FetchActiveTasks(db, limit)
		case false:
			tasks, err = pers.FetchInActiveTasks(db, limit)
		}
		return tasksFetched{tasks, active, err}
	}
}

func openTextEditor(fPath string, editorCmd []string, taskIndex int, taskID uint64, oldContext *string) tea.Cmd {
	c := exec.Command(editorCmd[0], append(editorCmd[1:], fPath)...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return tea.Msg(textEditorClosed{fPath, taskIndex, taskID, oldContext, err})
	})
}

func openURI(uri string) tea.Cmd {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	c := exec.Command(cmd, append(args, uri)...)
	err := c.Run()
	return func() tea.Msg {
		return uriOpenedMsg{uri, err}
	}
}

func openURIsDarwin(uris []string) tea.Cmd {
	c := exec.Command("open", uris...)
	err := c.Run()
	return func() tea.Msg {
		return urisOpenedDarwinMsg{uris, err}
	}
}

func copyContextToClipboard(context string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(context)
		return contextWrittenToCBMsg{err}
	}
}
