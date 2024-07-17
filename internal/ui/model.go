package ui

import (
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/utils"
)

const (
	compactListHeight = 10
	prefixPadding     = 20
	timeFormat        = "2006/01/02 15:04"
)

type itemDelegate struct {
	selStyle lipgloss.Style
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(types.Task)
	if !ok {
		return
	}

	start, _ := m.Paginator.GetSliceBounds(index)
	si := (index - start) % m.Paginator.PerPage

	var summ string
	summEls := strings.Split(t.Summary, ":")
	if len(summEls) > 1 {
		prefix := utils.RightPadTrim(summEls[0], prefixPadding, true)
		summ = prefix + strings.Join(summEls[1:], ":")
	} else {
		summ = t.Summary
	}
	var hasContext string
	if t.Context != nil {
		hasContext = "(c)"
	}
	str := fmt.Sprintf("[%d]\t%s%s", si+1, utils.RightPadTrim(summ, TaskSummaryMaxLen, true), hasContext)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.selStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

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
	helpView
)

type taskListType uint

const (
	activeTasks taskListType = iota
	archivedTasks
)

type model struct {
	db                 *sql.DB
	cfg                Config
	taskList           list.Model
	archivedTaskList   list.Model
	contextBMList      list.Model
	taskIndex          int
	taskId             uint64
	taskChange         taskChangeType
	contextVP          viewport.Model
	contextVPReady     bool
	taskDetailsVP      viewport.Model
	taskDetailsVPReady bool
	helpVP             viewport.Model
	helpVPReady        bool
	quitting           bool
	showHelpIndicator  bool
	errorMsg           string
	taskInput          textinput.Model
	activeView         activeView
	lastActiveView     activeView
	lastActiveList     taskListType
	tlTitleStyle       lipgloss.Style
	atlTitleStyle      lipgloss.Style
	tlSelStyle         lipgloss.Style
	atlSelStyle        lipgloss.Style
	terminalWidth      int
	terminalHeight     int
	contextVPTaskId    uint64
	rtos               string
	urlRegex           *regexp.Regexp
}
