package ui

import (
	"database/sql"
	"runtime"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/ui/theme"
	"github.com/dhth/omm/internal/utils"
)

func InitialModel(db *sql.DB, config Config, thm theme.Theme) Model {
	styles := newStyles(thm)

	taskItems := make([]list.Item, 0)

	taskList := list.New(taskItems,
		newTaskListDelegate(thm, config.ListDensity, activeTasks),
		taskSummaryWidth,
		defaultListHeight,
	)
	taskList.Title = config.TaskListTitle
	taskList.SetFilteringEnabled(true)
	taskList.SetStatusBarItemName("task", "tasks")
	taskList.SetShowStatusBar(true)
	taskList.SetShowHelp(false)
	taskList.DisableQuitKeybindings()
	taskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	taskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	taskList.SetStatusBarItemName("task", "tasks")

	taskList.Styles.Title = styles.activeListTitleBar

	archivedTaskItems := make([]list.Item, 0)

	archivedTaskList := list.New(archivedTaskItems,
		newTaskListDelegate(thm, config.ListDensity, archivedTasks),
		taskSummaryWidth,
		defaultListHeight,
	)
	archivedTaskList.Title = archivedTitle
	archivedTaskList.SetShowStatusBar(true)
	archivedTaskList.SetStatusBarItemName("task", "tasks")
	archivedTaskList.SetFilteringEnabled(true)
	archivedTaskList.SetShowHelp(false)
	archivedTaskList.DisableQuitKeybindings()
	archivedTaskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	archivedTaskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	archivedTaskList.SetStatusBarItemName("task", "tasks")

	archivedTaskList.Styles.Title = styles.archivedListTitleBar

	taskInput := textinput.New()
	taskInput.Placeholder = "prefix: task summary goes here"
	taskInput.CharLimit = types.TaskSummaryMaxLen
	taskInput.SetWidth(taskSummaryWidth)

	contextBMList := list.New(nil, newBookmarksListDelegate(thm), taskSummaryWidth, defaultListHeight)

	contextBMList.Title = "task bookmarks"
	contextBMList.SetShowHelp(false)
	contextBMList.SetStatusBarItemName("bookmark", "bookmarks")
	contextBMList.SetFilteringEnabled(false)
	contextBMList.DisableQuitKeybindings()
	contextBMList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	contextBMList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	contextBMList.Styles.Title = styles.bookmarksListTitleBar

	prefixSearchList := list.New(nil, newPrefixSearchListDelegate(thm), taskSummaryWidth, defaultListHeight)

	prefixSearchList.Title = "filter by prefix"
	prefixSearchList.SetShowHelp(false)
	prefixSearchList.SetStatusBarItemName("prefix", "prefixes")
	prefixSearchList.SetFilteringEnabled(false)
	prefixSearchList.DisableQuitKeybindings()
	prefixSearchList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	prefixSearchList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	prefixSearchList.Styles.Title = styles.prefixListTitleBar

	m := Model{
		db:                db,
		cfg:               config,
		theme:             thm,
		styles:            styles,
		taskList:          taskList,
		archivedTaskList:  archivedTaskList,
		taskBMList:        contextBMList,
		prefixSearchList:  prefixSearchList,
		taskInput:         taskInput,
		showHelpIndicator: true,
		contextVPTaskID:   0,
		rtos:              runtime.GOOS,
		uriRegex:          utils.GetURIRegex(),
	}

	return m
}
