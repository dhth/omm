package ui

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
)

const (
	defaultListHeight = 10
	prefixPadding     = 20
	timeFormat        = "2006/01/02 15:04"
	taskSummaryWidth  = 120
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchTasks(m.db, true, pers.TaskNumLimit),
		fetchTasks(m.db, false, pers.TaskNumLimit),
		hideHelp(time.Minute*1),
	)
}

type taskChangeType uint

const (
	taskInsert taskChangeType = iota
	taskUpdateSummary
	taskChangePriority
)

type activeView uint

const (
	taskListView activeView = iota
	archivedTaskListView
	taskEntryView
	taskDetailsView
	contextBookmarksView
	prefixSearchView
	helpView
)

type taskListType uint

const (
	activeTasks taskListType = iota
	archivedTasks
)

type model struct {
	db                    *sql.DB
	cfg                   Config
	taskList              list.Model
	archivedTaskList      list.Model
	taskBMList            list.Model
	prefixSearchList      list.Model
	activeTasksPrefixes   map[types.TaskPrefix]struct{}
	archivedTasksPrefixes map[types.TaskPrefix]struct{}
	tlIndexMap            map[uint64]int
	taskIndex             int
	taskId                uint64
	taskChange            taskChangeType
	contextVP             viewport.Model
	contextVPReady        bool
	taskDetailsVP         viewport.Model
	taskDetailsVPReady    bool
	helpVP                viewport.Model
	helpVPReady           bool
	quitting              bool
	showHelpIndicator     bool
	successMsg            string
	errorMsg              string
	taskInput             textinput.Model
	activeView            activeView
	lastActiveView        activeView
	activeTaskList        taskListType
	tlTitleStyle          lipgloss.Style
	atlTitleStyle         lipgloss.Style
	tlSelStyle            lipgloss.Style
	atlSelStyle           lipgloss.Style
	terminalWidth         int
	terminalHeight        int
	contextVPTaskId       uint64
	rtos                  string
	urlRegex              *regexp.Regexp
	shortenedListHt       int
	contextMdRenderer     *glamour.TermRenderer
	taskDetailsMdRenderer *glamour.TermRenderer
}
