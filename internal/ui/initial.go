package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	TaskSummaryMaxLen = 100
)

func InitialModel(db *sql.DB, config Config) model {

	taskItems := make([]list.Item, 0)
	tlSelItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.TaskListColor))

	var taskList list.Model
	switch config.ListDensity {
	case Compact:
		taskList = list.New(taskItems, itemDelegate{selStyle: tlSelItemStyle}, TaskSummaryMaxLen, compactListHeight)
		taskList.SetShowStatusBar(false)
	case Spacious:
		taskList = list.New(taskItems, newItemDelegate(lipgloss.Color(config.TaskListColor)), TaskSummaryMaxLen, 14)
		taskList.SetShowStatusBar(true)
	}
	taskList.SetShowTitle(false)
	taskList.SetFilteringEnabled(false)
	taskList.SetShowHelp(false)
	taskList.DisableQuitKeybindings()
	taskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	taskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	taskList.Styles.Title = taskList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(config.TaskListColor)).
		Bold(true)
	taskListTitleStyle := titleStyle.Background(lipgloss.Color(config.TaskListColor))

	atlSelItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.ArchivedTaskListColor))
	archivedTaskItems := make([]list.Item, 0)

	var archivedTaskList list.Model
	switch config.ListDensity {
	case Compact:
		archivedTaskList = list.New(archivedTaskItems, itemDelegate{selStyle: atlSelItemStyle}, TaskSummaryMaxLen, compactListHeight)
		archivedTaskList.SetShowStatusBar(false)
	case Spacious:
		archivedTaskList = list.New(archivedTaskItems, newItemDelegate(lipgloss.Color(config.ArchivedTaskListColor)), TaskSummaryMaxLen, 16)
		archivedTaskList.SetShowStatusBar(true)
	}
	archivedTaskList.SetShowTitle(false)
	archivedTaskList.SetFilteringEnabled(false)
	archivedTaskList.SetShowHelp(false)
	archivedTaskList.DisableQuitKeybindings()
	archivedTaskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	archivedTaskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	archivedTaskList.Styles.Title = archivedTaskList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(config.ArchivedTaskListColor)).
		Bold(true)
	archivedTaskListTitleStyle := titleStyle.Background(lipgloss.Color(config.ArchivedTaskListColor))

	taskInput := textinput.New()
	taskInput.Placeholder = "prefix: task summary goes here"
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
		tlSelStyle:        tlSelItemStyle,
		atlSelStyle:       atlSelItemStyle,
		contextVPTaskId:   0,
	}

	return m
}
