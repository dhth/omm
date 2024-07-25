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
	var spacing int
	if config.ShowContext {
		spacing = 1
	}
	switch config.ListDensity {
	case Compact:
		taskList = list.New(taskItems, newListDelegate(lipgloss.Color(config.TaskListColor), false, spacing), taskSummaryWidth, defaultListHeight)
	case Spacious:
		taskList = list.New(taskItems, newListDelegate(lipgloss.Color(config.TaskListColor), true, 1), taskSummaryWidth, defaultListHeight)
	}
	taskList.Title = config.TaskListTitle
	taskList.SetFilteringEnabled(true)
	taskList.SetStatusBarItemName("task", "tasks")
	taskList.SetShowStatusBar(true)
	taskList.SetShowHelp(false)
	taskList.DisableQuitKeybindings()
	taskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	taskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	taskList.SetStatusBarItemName("task", "tasks")

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
		archivedTaskList = list.New(archivedTaskItems, newListDelegate(lipgloss.Color(config.ArchivedTaskListColor), false, spacing), taskSummaryWidth, defaultListHeight)
	case Spacious:
		archivedTaskList = list.New(archivedTaskItems, newListDelegate(lipgloss.Color(config.ArchivedTaskListColor), true, 1), taskSummaryWidth, defaultListHeight)
	}
	archivedTaskList.Title = "archived"
	archivedTaskList.SetShowStatusBar(true)
	archivedTaskList.SetStatusBarItemName("task", "tasks")
	archivedTaskList.SetFilteringEnabled(true)
	archivedTaskList.SetShowHelp(false)
	archivedTaskList.DisableQuitKeybindings()
	archivedTaskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	archivedTaskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	archivedTaskList.SetStatusBarItemName("task", "tasks")

	archivedTaskList.Styles.Title = archivedTaskList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(config.ArchivedTaskListColor)).
		Bold(true)
	archivedTaskListTitleStyle := titleStyle.Background(lipgloss.Color(config.ArchivedTaskListColor))

	taskInput := textinput.New()
	taskInput.Placeholder = "prefix: task summary goes here"
	taskInput.CharLimit = types.TaskSummaryMaxLen
	taskInput.Width = taskSummaryWidth

	contextBMList := list.New(nil, newListDelegate(lipgloss.Color(contextBMColor), false, 1), taskSummaryWidth, defaultListHeight)

	contextBMList.Title = "task bookmarks"
	contextBMList.SetShowHelp(false)
	contextBMList.SetStatusBarItemName("bookmark", "bookmarks")
	contextBMList.SetFilteringEnabled(false)
	contextBMList.DisableQuitKeybindings()
	contextBMList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	contextBMList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	contextBMList.Styles.Title = contextBMList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(contextBMColor)).
		Bold(true)

	prefixSearchList := list.New(nil, newListDelegate(lipgloss.Color(prefixSearchColor), false, 0), taskSummaryWidth, defaultListHeight)

	prefixSearchList.Title = "filter by prefix"
	prefixSearchList.SetShowHelp(false)
	prefixSearchList.SetStatusBarItemName("prefix", "prefixes")
	prefixSearchList.SetFilteringEnabled(false)
	prefixSearchList.DisableQuitKeybindings()
	prefixSearchList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	prefixSearchList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	prefixSearchList.Styles.Title = prefixSearchList.Styles.Title.
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color(prefixSearchColor)).
		Bold(true)

	activeTasksPrefixes := make(map[types.TaskPrefix]struct{})
	archivedTasksPrefixes := make(map[types.TaskPrefix]struct{})

	tr, _ := getMarkDownRenderer(taskDetailsWordWrap)

	m := model{
		db:                    db,
		cfg:                   config,
		taskList:              taskList,
		archivedTaskList:      archivedTaskList,
		taskBMList:            contextBMList,
		prefixSearchList:      prefixSearchList,
		activeTasksPrefixes:   activeTasksPrefixes,
		archivedTasksPrefixes: archivedTasksPrefixes,
		taskInput:             taskInput,
		showHelpIndicator:     true,
		tlTitleStyle:          taskListTitleStyle,
		atlTitleStyle:         archivedTaskListTitleStyle,
		tlSelStyle:            tlSelItemStyle,
		atlSelStyle:           atlSelItemStyle,
		contextVPTaskId:       0,
		rtos:                  runtime.GOOS,
		urlRegex:              xurls.Strict(),
		taskDetailsMdRenderer: tr,
	}

	return m
}
