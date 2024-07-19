package ui

import (
	"database/sql"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/types"
	"mvdan.cc/xurls/v2"
)

func InitialModel(db *sql.DB, config Config) model {

	taskItems := make([]list.Item, 0)
	tlSelItemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(config.TaskListColor))

	var taskList list.Model
	switch config.ListDensity {
	case Compact:
		taskList = list.New(taskItems, itemDelegate{selStyle: tlSelItemStyle}, taskSummaryWidth, compactListHeight)
		taskList.SetShowStatusBar(false)
	case Spacious:
		taskList = list.New(taskItems, newTaskListDelegate(lipgloss.Color(config.TaskListColor)), taskSummaryWidth, 14)
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
		archivedTaskList = list.New(archivedTaskItems, itemDelegate{selStyle: atlSelItemStyle}, taskSummaryWidth, compactListHeight)
		archivedTaskList.SetShowStatusBar(false)
	case Spacious:
		archivedTaskList = list.New(archivedTaskItems, newTaskListDelegate(lipgloss.Color(config.ArchivedTaskListColor)), taskSummaryWidth, 16)
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
	taskInput.CharLimit = types.TaskSummaryMaxLen
	taskInput.Width = taskSummaryWidth

	contextBMList := list.New(nil, newContextURLListDel(contextBMColor), taskSummaryWidth, compactListHeight)

	contextBMList.SetShowTitle(false)
	contextBMList.SetShowHelp(false)
	contextBMList.SetFilteringEnabled(false)
	contextBMList.DisableQuitKeybindings()
	contextBMList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	contextBMList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	m := model{
		db:                db,
		cfg:               config,
		taskList:          taskList,
		archivedTaskList:  archivedTaskList,
		contextBMList:     contextBMList,
		taskInput:         taskInput,
		showHelpIndicator: true,
		tlTitleStyle:      taskListTitleStyle,
		atlTitleStyle:     archivedTaskListTitleStyle,
		tlSelStyle:        tlSelItemStyle,
		atlSelStyle:       atlSelItemStyle,
		contextVPTaskId:   0,
		rtos:              runtime.GOOS,
		urlRegex:          xurls.Strict(),
	}

	return m
}
