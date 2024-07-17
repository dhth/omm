package ui

import (
	"os/exec"
	"runtime"
	"time"

	"database/sql"

	"github.com/atotto/clipboard"
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
		now := time.Now()
		err := pers.UpdateTaskSummaryInDB(db, id, summary, now)
		return taskSummaryUpdatedMsg{listIndex, id, summary, now, err}
	}
}

func updateTaskContext(db *sql.DB, listIndex int, id uint64, context string, list taskListType) tea.Cmd {
	return func() tea.Msg {
		var err error
		now := time.Now()
		if context == "" {
			err = pers.UnsetTaskContextInDB(db, id, now)
		} else {
			err = pers.UpdateTaskContextInDB(db, id, context, now)
		}
		return taskContextUpdatedMsg{listIndex, list, id, context, now, err}
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

func openTextEditor(fPath string, editorCmd []string, taskIndex int, taskId uint64, oldContext *string) tea.Cmd {

	c := exec.Command(editorCmd[0], append(editorCmd[1:], fPath)...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return tea.Msg(textEditorClosed{fPath, taskIndex, taskId, oldContext, err})
	})
}

func openURL(url string) tea.Cmd {
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
	c := exec.Command(cmd, append(args, url)...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return urlOpenedMsg{url, err}
	})
}

func openURLsDarwin(urls []string) tea.Cmd {
	c := exec.Command("open", urls...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return urlsOpenedDarwinMsg{urls, err}
	})
}

func copyContextToClipboard(context string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(context)
		return contextWrittenToCBMsg{err}
	}
}
