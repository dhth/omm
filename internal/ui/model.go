package ui

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/utils"
)

type itemDelegate struct {
	selStyle lipgloss.Style
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(types.Task)
	if !ok {
		return
	}

	start, _ := m.Paginator.GetSliceBounds(index)
	si := (index - start) % m.Paginator.PerPage

	str := fmt.Sprintf("[%d]\t%s", si+1, utils.RightPadTrim(i.Summary, 80, true))

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
		fetchTasks(m.db, true, 50),
		fetchTasks(m.db, false, 50),
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
	helpView
)

type model struct {
	db                *sql.DB
	cfg               Config
	taskList          list.Model
	archivedTaskList  list.Model
	taskIndex         int
	taskId            uint64
	taskChange        taskChangeType
	quitting          bool
	showHelpIndicator bool
	message           string
	taskInput         textinput.Model
	activeView        activeView
	tlTitleStyle      lipgloss.Style
	atlTitleStyle     lipgloss.Style
}
