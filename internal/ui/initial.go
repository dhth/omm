package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	TaskSummaryMaxLen = 80
)

func InitialModel(db *sql.DB, config Config) model {

	taskItems := make([]list.Item, 0)
	tlSelItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.TaskListColor))

	archivedTLSelItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.ArchivedTaskListColor))

	taskList := list.New(taskItems, itemDelegate{selStyle: tlSelItemStyle}, TaskSummaryMaxLen, 10)
	taskList.SetShowStatusBar(false)
	taskList.SetFilteringEnabled(false)
	taskList.SetShowHelp(false)
	taskList.SetShowTitle(false)
	taskList.DisableQuitKeybindings()
	taskList.Styles.Title = taskList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(TaskListColor)).
		Bold(true)

	taskListTitleStyle := titleStyle.Background(lipgloss.Color(config.TaskListColor))

	archivedTaskItems := make([]list.Item, 0)

	archivedTaskList := list.New(archivedTaskItems, itemDelegate{selStyle: archivedTLSelItemStyle}, TaskSummaryMaxLen, 10)
	archivedTaskList.SetShowStatusBar(false)
	archivedTaskList.SetFilteringEnabled(false)
	archivedTaskList.SetShowHelp(false)
	archivedTaskList.SetShowTitle(false)
	archivedTaskList.DisableQuitKeybindings()
	archivedTaskList.Styles.Title = archivedTaskList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(ArchivedTLColor)).
		Bold(true)
	archivedTaskListTitleStyle := titleStyle.Background(lipgloss.Color(config.ArchivedTaskListColor))

	taskInput := textinput.New()
	taskInput.Placeholder = "task summary goes here"
	taskInput.CharLimit = TaskSummaryMaxLen
	taskInput.Width = TaskSummaryMaxLen

	m := model{
		db:                db,
		cfg:               config,
		taskList:          taskList,
		archivedTaskList:  archivedTaskList,
		taskInput:         taskInput,
		showHelpIndicator: true,
		tlTitleStyle:      taskListTitleStyle,
		atlTitleStyle:     archivedTaskListTitleStyle,
	}

	return m
}
