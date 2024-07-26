package ui

import (
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/utils"
)

const (
	noSpaceAvailableMsg   = "Task list is at capacity. Archive/delete tasks using ctrl+d/ctrl+x."
	noContextMsg          = "  ∅"
	viewPortMoveLineCount = 3
)

//go:embed assets/help.md
var helpStr string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.successMsg = ""
	m.errorMsg = ""

	if m.activeView == taskListView || m.activeView == archivedTaskListView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if m.taskList.FilterState() == list.Filtering {
				m.taskList, cmd = m.taskList.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			if m.archivedTaskList.FilterState() == list.Filtering {
				m.archivedTaskList, cmd = m.archivedTaskList.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}
	}

	if m.activeView == taskEntryView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "esc", "ctrl+c":
				m.activeView = taskListView
				m.activeTaskList = activeTasks
			case "enter":
				taskSummary := m.taskInput.Value()
				taskSummary = strings.TrimSpace(taskSummary)

				if taskSummary == "" {
					m.activeView = taskListView
					m.activeTaskList = activeTasks
					break
				}

				summEls := strings.Split(taskSummary, types.PrefixDelimiter)
				if len(summEls) > 1 {
					if summEls[0] == "" {
						m.errorMsg = "prefix cannot be empty"
						break
					}
				}

				switch m.taskChange {
				case taskInsert:
					now := time.Now()
					cmd = createTask(m.db, taskSummary, now, now)
					cmds = append(cmds, cmd)
					m.taskInput.Reset()
					m.activeView = taskListView
					m.activeTaskList = activeTasks
				case taskUpdateSummary:
					cmd = updateTaskSummary(m.db, m.taskIndex, m.taskId, taskSummary)
					cmds = append(cmds, cmd)
					m.taskInput.Reset()
					m.activeView = taskListView
					m.activeTaskList = activeTasks
				}

			case "ctrl+p":
				if len(m.taskList.Items()) == 0 {
					m.errorMsg = "No items in task list"
					break
				}

				tasksPrefixes := make(map[types.TaskPrefix]struct{})

				summary := m.taskInput.Value()

				currentPrefix, prefixPresent := getPrefix(summary)

				for _, li := range m.taskList.Items() {
					t, ok := li.(types.Task)
					if !ok {
						continue
					}

					prefix, pOk := t.Prefix()
					if !pOk {
						continue
					}

					if prefixPresent && prefix.FilterValue() == currentPrefix {
						continue
					}

					tasksPrefixes[prefix] = struct{}{}
				}

				var prefixes []types.TaskPrefix
				for k := range tasksPrefixes {
					prefixes = append(prefixes, k)
				}

				if len(prefixes) == 0 {
					m.errorMsg = "No prefixes in task list"
					break
				}

				if len(prefixes) == 1 {
					m.errorMsg = "Only 1 unique prefix in task list"
					break
				}

				sort.Slice(prefixes, func(i, j int) bool {
					return prefixes[i] < prefixes[j]
				})

				pi := make([]list.Item, len(prefixes))
				for i, p := range prefixes {
					pi[i] = list.Item(p)
				}

				m.prefixSearchList.SetItems(pi)
				m.lastActiveView = m.activeView
				m.activeView = prefixSelectionView
				switch prefixPresent {
				case true:
					m.prefixSearchList.Title = "change prefix"
				case false:
					m.prefixSearchList.Title = "choose prefix"
				}
				m.prefixSearchUse = prefixChoose

				return m, tea.Batch(cmds...)
			}
		}

		m.taskInput, cmd = m.taskInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := listStyle.GetFrameSize()
		_, h3 := statusBarMsgStyle.GetFrameSize()
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		m.taskList.SetWidth(msg.Width - w)
		m.archivedTaskList.SetWidth(msg.Width - 2)
		m.taskBMList.SetWidth(msg.Width - 2)
		m.taskBMList.SetHeight(msg.Height - h - h3 - 1)
		m.prefixSearchList.SetWidth(msg.Width - 2)
		m.prefixSearchList.SetHeight(msg.Height - h - h3 - 1)

		var listHeight int
		contextHeight := (msg.Height - h - h3 - 5) / 2

		m.shortenedListHt = msg.Height - contextHeight - 5

		if m.cfg.ShowContext {
			listHeight = m.shortenedListHt
		} else {
			listHeight = msg.Height - h - h3 - 1
		}

		m.taskList.SetHeight(listHeight)
		m.archivedTaskList.SetHeight(listHeight)

		if !m.contextVPReady {
			m.contextVP = viewport.New(msg.Width-3, contextHeight)
			m.contextVPReady = true
		} else {
			m.contextVP.Width = msg.Width - 3
			m.contextVP.Height = contextHeight
		}

		if !m.taskDetailsVPReady {
			m.taskDetailsVP = viewport.New(msg.Width-4, m.terminalHeight-4)
			m.taskDetailsVP.KeyMap.HalfPageDown.SetKeys("ctrl+d")
			m.taskDetailsVPReady = true
			m.taskDetailsVP.KeyMap.Up.SetEnabled(false)
			m.taskDetailsVP.KeyMap.Down.SetEnabled(false)
		} else {
			m.taskDetailsVP.Width = msg.Width - 4
			m.taskDetailsVP.Height = m.terminalHeight - 4
		}

		crWrap := (msg.Width - 4)
		if crWrap > contextWordWrapUpperLimit {
			crWrap = contextWordWrapUpperLimit
		}
		m.contextMdRenderer, _ = utils.GetMarkDownRenderer(crWrap)

		helpToRender := helpStr
		switch m.contextMdRenderer {
		case nil:
			break
		default:
			helpStrGl, err := m.contextMdRenderer.Render(helpStr)
			if err != nil {
				break
			}
			helpToRender = helpStrGl
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width-3, m.terminalHeight-4)
			m.helpVP.SetContent(helpToRender)
			m.helpVP.KeyMap.Up.SetEnabled(false)
			m.helpVP.KeyMap.Down.SetEnabled(false)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width - 3
			m.helpVP.Height = m.terminalHeight - 4
		}

	case tea.KeyMsg:
		if m.cfg.ConfirmBeforeDeletion && m.showDeletePrompt && msg.String() != "ctrl+x" {
			m.showDeletePrompt = false

			switch m.activeView {
			case taskListView:
				m.taskList.Title = m.cfg.TaskListTitle
				m.taskList.Styles.Title = m.taskList.Styles.Title.Background(lipgloss.Color(m.cfg.TaskListColor))
			case archivedTaskListView:
				m.archivedTaskList.Title = "archived"
				m.archivedTaskList.Styles.Title = m.archivedTaskList.Styles.Title.Background(lipgloss.Color(m.cfg.ArchivedTaskListColor))
			}
			return m, tea.Batch(cmds...)
		}

		switch keypress := msg.String(); keypress {

		case "Q":
			m.quitting = true
			if m.cfg.Guide {
				_ = os.Remove(m.cfg.DBPath)
			}
			return m, tea.Quit

		case "esc", "q", "ctrl+c":
			av := m.activeView

			if m.activeView == taskListView && m.taskList.IsFiltered() {
				m.taskList.ResetFilter()
				break
			}

			if m.activeView == archivedTaskListView && m.archivedTaskList.IsFiltered() {
				m.archivedTaskList.ResetFilter()
				break
			}

			if m.activeView == archivedTaskListView {
				m.activeView = taskListView
				m.activeTaskList = activeTasks
				m.lastActiveView = av
				break
			}

			if m.activeView == taskDetailsView || m.activeView == contextBookmarksView || m.activeView == helpView {
				m.activeView = m.lastActiveView
				switch m.activeView {
				case taskListView:
					m.activeTaskList = activeTasks
				case archivedTaskListView:
					m.activeTaskList = archivedTasks
				}
				break
			}

			if m.activeView == prefixSelectionView {
				m.activeView = m.lastActiveView
				if m.prefixSearchUse == prefixChoose {
					return m, tea.Batch(cmds...)
				}
				break
			}

			m.quitting = true
			if m.cfg.Guide {
				_ = os.Remove(m.cfg.DBPath)
			}
			return m, tea.Quit

		case "?":
			if m.activeView == taskDetailsView || m.activeView == contextBookmarksView || m.activeView == prefixSelectionView {
				break
			}

			if m.activeView == helpView {
				m.activeView = m.lastActiveView
				break
			}
			m.lastActiveView = m.activeView
			m.activeView = helpView

		case "tab", "shift+tab":
			switch m.activeView {
			case taskListView:
				m.activeView = archivedTaskListView
				m.activeTaskList = archivedTasks
				m.lastActiveView = m.activeView
			case archivedTaskListView:
				m.activeView = taskListView
				m.activeTaskList = activeTasks
				m.lastActiveView = m.activeView
			}

		case "I":
			if m.activeView != taskListView {
				break
			}

			if !m.isSpaceAvailable() {
				m.errorMsg = noSpaceAvailableMsg
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot add items when the task list is filtered"
				break
			}

			m.taskIndex = 0
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "O":
			if m.activeView != taskListView {
				break
			}

			if !m.isSpaceAvailable() {
				m.errorMsg = noSpaceAvailableMsg
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot add items when the task list is filtered"
				break
			}

			m.taskIndex = m.taskList.Index()
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "a", "o":
			if m.activeView != taskListView {
				break
			}

			if !m.isSpaceAvailable() {
				m.errorMsg = noSpaceAvailableMsg
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot add items when the task list is filtered"
				break
			}

			if len(m.taskList.Items()) == 0 {
				m.taskIndex = 0
			} else {
				m.taskIndex = m.taskList.Index() + 1
			}
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "A":
			if m.activeView != taskListView {
				break
			}

			if !m.isSpaceAvailable() {
				m.errorMsg = noSpaceAvailableMsg
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot add items when the task list is filtered"
				break
			}

			m.taskIndex = len(m.taskList.Items())
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "down", "j":
			if m.activeView != taskDetailsView && m.activeView != helpView {
				break
			}

			switch m.activeView {
			case taskDetailsView:
				if m.taskDetailsVP.AtBottom() {
					break
				}
				m.taskDetailsVP.LineDown(viewPortMoveLineCount)
			case helpView:
				if m.helpVP.AtBottom() {
					break
				}
				m.helpVP.LineDown(viewPortMoveLineCount)
			}

		case "up", "k":
			if m.activeView != taskDetailsView && m.activeView != helpView {
				break
			}

			switch m.activeView {
			case taskDetailsView:
				if m.taskDetailsVP.AtTop() {
					break
				}
				m.taskDetailsVP.LineUp(viewPortMoveLineCount)
			case helpView:
				if m.helpVP.AtTop() {
					break
				}
				m.helpVP.LineUp(viewPortMoveLineCount)
			}

		case "J":
			if m.activeView != taskListView {
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot move items when the task list is filtered"
				break
			}

			if len(m.taskList.Items()) == 0 {
				break
			}

			ci := m.taskList.Index()
			if ci == len(m.taskList.Items())-1 {
				break
			}

			itemAbove := m.taskList.Items()[ci+1]
			currentItem := m.taskList.Items()[ci]
			m.taskList.SetItem(ci, itemAbove)
			m.taskList.SetItem(ci+1, currentItem)
			m.taskList.Select(ci + 1)

			cmd = m.updateActiveTasksSequence()
			cmds = append(cmds, cmd)

		case "K":
			if m.activeView != taskListView {
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot move items when the task list is filtered"
				break
			}

			ci := m.taskList.Index()
			if ci == 0 {
				break
			}

			itemAbove := m.taskList.Items()[ci-1]
			currentItem := m.taskList.Items()[ci]
			m.taskList.SetItem(ci, itemAbove)
			m.taskList.SetItem(ci-1, currentItem)
			m.taskList.Select(ci - 1)

			cmd = m.updateActiveTasksSequence()
			cmds = append(cmds, cmd)

		case "u":
			if m.activeView != taskListView {
				break
			}

			if len(m.taskList.Items()) == 0 {
				break
			}

			listItem := m.taskList.SelectedItem()
			index := m.taskList.Index()
			t, ok := listItem.(types.Task)
			if !ok {
				m.errorMsg = "Something went wrong"
				break
			}

			m.taskInput.SetValue(t.Summary)
			m.taskInput.Focus()
			m.taskIndex = index
			m.taskId = t.ID
			m.taskChange = taskUpdateSummary
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "ctrl+r":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			br := false
			switch m.activeView {
			case taskListView:
				if m.taskList.IsFiltered() {
					br = true
				}

			case archivedTaskListView:
				if m.archivedTaskList.IsFiltered() {
					br = true
				}
			}

			if br {
				break
			}

			cmds = append(cmds, fetchTasks(m.db, true, pers.TaskNumLimit))
			cmds = append(cmds, fetchTasks(m.db, false, pers.TaskNumLimit))

		case "ctrl+d":
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				if m.taskList.IsFiltered() {
					m.errorMsg = "Cannot archive items when the task list is filtered"
					break
				}

				listItem := m.taskList.SelectedItem()
				index := m.taskList.Index()
				t, ok := listItem.(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong; cannot archive item"
					break
				}

				cmd = changeTaskStatus(m.db, index, t.ID, false, time.Now())
				cmds = append(cmds, cmd)

			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					break
				}

				if m.archivedTaskList.IsFiltered() {
					m.errorMsg = "Cannot unarchive items when the task list is filtered"
					break
				}

				listItem := m.archivedTaskList.SelectedItem()
				index := m.archivedTaskList.Index()
				t, ok := listItem.(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}

				cmd = changeTaskStatus(m.db, index, t.ID, true, time.Now())
				cmds = append(cmds, cmd)
			}

		case "ctrl+x":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			quit := false
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					quit = true
					break
				}

				if m.taskList.IsFiltered() {
					m.errorMsg = "Cannot delete items when the task list is filtered"
					quit = true
					break
				}
			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					quit = true
					break
				}

				if m.archivedTaskList.IsFiltered() {
					m.errorMsg = "Cannot delete items when the task list is filtered"
					quit = true
					break
				}
			}

			if quit {
				break
			}

			if m.cfg.ConfirmBeforeDeletion && !m.showDeletePrompt {
				m.showDeletePrompt = true

				switch m.activeView {
				case taskListView:
					m.taskList.Title = "delete ?"
					m.taskList.Styles.Title = m.taskList.Styles.Title.Background(lipgloss.Color(promptColor))
				case archivedTaskListView:
					m.archivedTaskList.Title = "delete ?"
					m.archivedTaskList.Styles.Title = m.archivedTaskList.Styles.Title.Background(lipgloss.Color(promptColor))
				}

				break
			}

			switch m.activeView {
			case taskListView:
				index := m.taskList.Index()
				t, ok := m.taskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}
				cmd = deleteTask(m.db, t.ID, index, true)
				cmds = append(cmds, cmd)
				if m.cfg.ConfirmBeforeDeletion {
					m.showDeletePrompt = false
					m.taskList.Title = m.cfg.TaskListTitle
					m.taskList.Styles.Title = m.taskList.Styles.Title.Background(lipgloss.Color(m.cfg.TaskListColor))
				}

			case archivedTaskListView:
				index := m.archivedTaskList.Index()
				task, ok := m.archivedTaskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}

				cmd = deleteTask(m.db, task.ID, index, false)
				cmds = append(cmds, cmd)
				if m.cfg.ConfirmBeforeDeletion {
					m.showDeletePrompt = false
					m.archivedTaskList.Title = "archived"
					m.archivedTaskList.Styles.Title = m.archivedTaskList.Styles.Title.Background(lipgloss.Color(m.cfg.ArchivedTaskListColor))
				}
			}

		case "ctrl+p":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			var taskList list.Model
			taskPrefixes := make(map[types.TaskPrefix]struct{})

			switch m.activeView {
			case taskListView:
				taskList = m.taskList
			case archivedTaskListView:
				taskList = m.archivedTaskList
			}

			if len(taskList.Items()) == 0 {
				m.errorMsg = "No items in task list"
				break
			}

			for _, li := range taskList.Items() {
				t, ok := li.(types.Task)
				if ok {
					prefix, pOk := t.Prefix()
					if pOk {
						taskPrefixes[prefix] = struct{}{}
					}
				}
			}
			var prefixes []types.TaskPrefix
			for k := range taskPrefixes {
				prefixes = append(prefixes, k)
			}

			if len(prefixes) == 0 {
				m.errorMsg = "No prefixes in task list"
				break
			}

			if len(prefixes) == 1 {
				m.errorMsg = "Only 1 unique prefix in task list"
				break
			}

			sort.Slice(prefixes, func(i, j int) bool {
				return prefixes[i] < prefixes[j]
			})

			pi := make([]list.Item, len(prefixes))
			for i, p := range prefixes {
				pi[i] = list.Item(p)
			}

			m.prefixSearchList.SetItems(pi)
			m.lastActiveView = m.activeView
			m.activeView = prefixSelectionView
			m.prefixSearchList.Title = "filter by prefix"
			m.prefixSearchUse = prefixFilter

		case "enter":
			if m.activeView != taskListView && m.activeView != archivedTaskListView && m.activeView != contextBookmarksView && m.activeView != prefixSelectionView {
				break
			}
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				if m.taskList.IsFiltered() {
					selected, ok := m.taskList.SelectedItem().(types.Task)
					if !ok {
						m.errorMsg = "Something went wrong"
						break
					}

					listIndex, ok := m.tlIndexMap[selected.ID]
					if !ok {
						m.errorMsg = "Something went wrong"
						break
					}

					m.taskList.ResetFilter()
					m.taskList.Select(listIndex)
					break
				}

				index := m.taskList.Index()

				if index == 0 {
					m.errorMsg = "This item is already at the top of the list"
					break
				}

				listItem := m.taskList.SelectedItem()
				m.taskList.RemoveItem(index)
				cmd = m.taskList.InsertItem(0, listItem)
				cmds = append(cmds, cmd)
				m.taskList.Select(0)

				cmd = m.updateActiveTasksSequence()
				cmds = append(cmds, cmd)

			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					break
				}

				if !m.archivedTaskList.IsFiltered() {
					break
				}

				selected, ok := m.archivedTaskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}

				listIndex, ok := m.atlIndexMap[selected.ID]
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}

				m.archivedTaskList.ResetFilter()
				m.archivedTaskList.Select(listIndex)

			case contextBookmarksView:
				url := m.taskBMList.SelectedItem().FilterValue()
				cmds = append(cmds, openURL(url))
			case prefixSelectionView:
				prefix := m.prefixSearchList.SelectedItem().FilterValue()

				switch m.prefixSearchUse {
				case prefixFilter:
					var taskList list.Model

					switch m.activeTaskList {
					case activeTasks:
						taskList = m.taskList
					case archivedTasks:
						taskList = m.archivedTaskList
					}

					taskList.ResetFilter()
					var tlCmd tea.Cmd

					runes := []rune(prefix)

					if len(runes) > 1 {
						taskList.FilterInput.SetValue(string(runes[:len(runes)-1]))
					}

					taskList, tlCmd = taskList.Update(tea.KeyMsg{Type: -1, Runes: []int32{47}, Alt: false, Paste: false})
					cmds = append(cmds, tlCmd)

					taskList, tlCmd = taskList.Update(tea.KeyMsg{Type: -1, Runes: []rune{runes[len(runes)-1]}, Alt: false, Paste: false})
					cmds = append(cmds, tlCmd)

					// TODO: Try sending ENTER programmatically too
					// taskList, tlCmd = taskList.Update(tea.KeyMsg{Type: 13, Runes: []int32(nil), Alt: false, Paste: false})
					// or
					// taskList, tlCmd = taskList.Update(tea.KeyEnter)
					// this results in the list's paginator being broken, so requires another manual ENTER keypress

					switch m.activeTaskList {
					case activeTasks:
						m.taskList = taskList
						m.activeView = taskListView
					case archivedTasks:
						m.archivedTaskList = taskList
						m.activeView = archivedTaskListView
					}

					return m, tea.Sequence(cmds...)

				case prefixChoose:
					m.taskInput.SetValue(getSummaryWithNewPrefix(m.taskInput.Value(), prefix))
					m.activeView = taskEntryView
				}

			}

		case "E":
			if m.activeView != taskListView {
				break
			}

			if len(m.taskList.Items()) == 0 {
				break
			}

			if m.taskList.IsFiltered() {
				m.errorMsg = "Cannot move items when the task list is filtered"
				break
			}

			index := m.taskList.Index()

			lastIndex := len(m.taskList.Items()) - 1

			if index == lastIndex {
				m.errorMsg = "This item is already at the end of the list"
				break
			}

			if m.taskList.IsFiltered() {
				m.taskList.ResetFilter()
				m.taskList.Select(index)
			}

			listItem := m.taskList.SelectedItem()
			m.taskList.RemoveItem(index)
			cmd = m.taskList.InsertItem(lastIndex, listItem)
			cmds = append(cmds, cmd)
			m.taskList.Select(lastIndex)

			cmd = m.updateActiveTasksSequence()
			cmds = append(cmds, cmd)

		case "c":
			if m.activeView != taskListView && m.activeView != archivedTaskListView && m.activeView != taskDetailsView {
				break
			}

			var t types.Task
			var ok bool
			var index int

			switch m.activeTaskList {
			case activeTasks:
				t, ok = m.taskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}
				index = m.taskList.Index()
			case archivedTasks:
				t, ok = m.archivedTaskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}
				index = m.archivedTaskList.Index()
			}

			if len(m.cfg.TextEditorCmd) == 0 {
				m.errorMsg = "No editor has been set via --editor, or $EDITOR or $VISUAL"
				break
			}

			tempFile, err := os.CreateTemp("", "omm-*.md")
			if err != nil {
				m.errorMsg = fmt.Sprintf("Error creating temporary file: %s", err)
				break
			}
			if t.Context != nil {
				_, err = tempFile.Write([]byte(*t.Context))
				if err != nil {
					_ = tempFile.Close()
					break
				}
			}
			_ = tempFile.Close()

			cmds = append(cmds, openTextEditor(tempFile.Name(), m.cfg.TextEditorCmd, index, t.ID, t.Context))

		case "v":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			var tlDel list.ItemDelegate
			var atlDel list.ItemDelegate

			switch m.cfg.ListDensity {
			case Compact:
				tlDel = newSpaciousListDelegate(lipgloss.Color(m.cfg.TaskListColor), true, 1)
				atlDel = newSpaciousListDelegate(lipgloss.Color(m.cfg.ArchivedTaskListColor), true, 1)

				m.cfg.ListDensity = Spacious

			case Spacious:
				tlDel = compactItemDelegate{m.tlSelStyle}
				atlDel = compactItemDelegate{m.atlSelStyle}
				m.cfg.ListDensity = Compact
			}

			m.taskList.SetDelegate(tlDel)
			m.archivedTaskList.SetDelegate(atlDel)

			if m.cfg.ShowContext {
				m.taskList.SetHeight(m.shortenedListHt)
				m.archivedTaskList.SetHeight(m.shortenedListHt)
			}

		case "C":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			m.cfg.ShowContext = !m.cfg.ShowContext

			_, h := listStyle.GetFrameSize()
			_, h3 := statusBarMsgStyle.GetFrameSize()
			var listHeight int

			if m.cfg.ShowContext {
				listHeight = m.shortenedListHt
			} else {
				listHeight = m.terminalHeight - h - h3 - 1
			}

			if m.cfg.ListDensity == Compact {
				tlDel := compactItemDelegate{m.tlSelStyle}
				atlDel := compactItemDelegate{m.atlSelStyle}
				m.taskList.SetDelegate(tlDel)
				m.archivedTaskList.SetDelegate(atlDel)
			}

			m.taskList.SetHeight(listHeight)
			m.archivedTaskList.SetHeight(listHeight)

		case "d":
			if m.activeView == taskDetailsView {
				m.activeView = m.lastActiveView
				break
			}

			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			var t types.Task
			var ok bool

			switch m.activeView {
			case taskListView:
				t, ok = m.taskList.SelectedItem().(types.Task)
			case archivedTaskListView:
				t, ok = m.archivedTaskList.SelectedItem().(types.Task)
			}

			if !ok {
				break
			}

			m.taskDetailsVP.GotoTop()
			m.setContextFSContent(t)

			switch m.activeView {
			case taskListView:
				m.activeTaskList = activeTasks
			default:
				m.activeTaskList = archivedTasks
			}
			m.lastActiveView = m.activeView
			m.activeView = taskDetailsView

		case "h":
			if m.activeView != taskDetailsView {
				break
			}
			var t types.Task
			var ok bool

			switch m.activeTaskList {
			case activeTasks:
				m.taskList.CursorUp()
				t, ok = m.taskList.SelectedItem().(types.Task)
			case archivedTasks:
				m.archivedTaskList.CursorUp()
				t, ok = m.archivedTaskList.SelectedItem().(types.Task)
			}

			if !ok {
				break
			}

			m.taskDetailsVP.GotoTop()
			m.setContextFSContent(t)

		case "l":
			if m.activeView != taskDetailsView {
				break
			}
			var t types.Task
			var ok bool

			switch m.activeTaskList {
			case activeTasks:
				m.taskList.CursorDown()
				t, ok = m.taskList.SelectedItem().(types.Task)
			case archivedTasks:
				m.archivedTaskList.CursorDown()
				t, ok = m.archivedTaskList.SelectedItem().(types.Task)
			}

			if !ok {
				break
			}

			m.taskDetailsVP.GotoTop()
			m.setContextFSContent(t)

		case "b":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			urls, ok := m.getTaskUrls()
			if !ok {
				break
			}

			if len(urls) == 0 {
				m.errorMsg = "No bookmarks for this task"
				break
			}

			if len(urls) == 1 {
				cmds = append(cmds, openURL(urls[0]))
				break
			}

			bmItems := make([]list.Item, len(urls))
			for i, url := range urls {
				bmItems[i] = list.Item(types.ContextBookmark(url))
			}
			m.taskBMList.SetItems(bmItems)
			switch m.activeView {
			case taskListView:
				m.activeTaskList = activeTasks
			case archivedTaskListView:
				m.activeTaskList = archivedTasks
			}
			m.lastActiveView = m.activeView
			m.activeView = contextBookmarksView

		case "B":
			if m.activeView != taskListView && m.activeView != archivedTaskListView && m.activeView != taskDetailsView {
				break
			}

			urls, ok := m.getTaskUrls()
			if !ok {
				break
			}

			if len(urls) == 0 {
				m.errorMsg = "No bookmarks for this task"
				break
			}

			if len(urls) == 1 {
				cmds = append(cmds, openURL(urls[0]))
				break
			}

			if m.rtos == types.GOOSDarwin {
				cmds = append(cmds, openURLsDarwin(urls))
				break
			}

			for _, url := range urls {
				cmds = append(cmds, openURL(url))
			}

		case "y":
			if m.activeView != taskListView && m.activeView != archivedTaskListView && m.activeView != taskDetailsView {
				break
			}

			var t types.Task
			var ok bool

			switch m.activeView {
			case taskListView:
				t, ok = m.taskList.SelectedItem().(types.Task)
			case archivedTaskListView:
				t, ok = m.archivedTaskList.SelectedItem().(types.Task)
			case taskDetailsView:
				switch m.activeTaskList {
				case activeTasks:
					t, ok = m.taskList.SelectedItem().(types.Task)
				case archivedTasks:
					t, ok = m.archivedTaskList.SelectedItem().(types.Task)
				}
			}

			if !ok {
				break
			}

			if t.Context == nil {
				m.errorMsg = "There's no context to copy"
				break
			}

			cmds = append(cmds, copyContextToClipboard(*t.Context))
		}

	case HideHelpMsg:
		m.showHelpIndicator = false

	case taskCreatedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error creating task: %s", msg.err)
			break
		}

		t := types.Task{
			ID:        msg.id,
			Summary:   msg.taskSummary,
			Active:    true,
			CreatedAt: msg.createdAt,
			UpdatedAt: msg.updatedAt,
		}
		entry := list.Item(t)
		cmd = m.taskList.InsertItem(m.taskIndex, entry)
		cmds = append(cmds, cmd)
		m.taskList.Select(m.taskIndex)

		cmd = m.updateActiveTasksSequence()
		cmds = append(cmds, cmd)

	case taskDeletedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error deleting task: %s", msg.err)
			break
		}

		switch msg.active {
		case true:
			m.taskList.RemoveItem(msg.listIndex)
			cmd = m.updateActiveTasksSequence()
			cmds = append(cmds, cmd)
		case false:
			m.archivedTaskList.RemoveItem(msg.listIndex)
			m.updateArchivedTasksIndex()
		}

	case taskSequenceUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error updating task sequence: %s", msg.err)
		}

	case taskSummaryUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error updating task: %s", msg.err)
		} else {
			listItem := m.taskList.Items()[msg.listIndex]
			t, ok := listItem.(types.Task)
			if !ok {
				break
			}

			t.Summary = msg.taskSummary
			t.UpdatedAt = msg.updatedAt
			cmd = m.taskList.SetItem(msg.listIndex, list.Item(t))
			cmds = append(cmds, cmd)
		}

	case taskContextUpdatedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error updating task: %s", msg.err)
		} else {
			var t types.Task
			var ok bool

			switch msg.list {
			case activeTasks:
				listItem := m.taskList.Items()[msg.listIndex]
				t, ok = listItem.(types.Task)
				if !ok {
					break
				}

				if msg.context == "" {
					t.Context = nil
				} else {
					t.Context = &msg.context
				}
				t.UpdatedAt = msg.updatedAt
				cmd = m.taskList.SetItem(msg.listIndex, list.Item(t))
				cmds = append(cmds, cmd)
			case archivedTasks:
				listItem := m.archivedTaskList.Items()[msg.listIndex]
				t, ok = listItem.(types.Task)
				if !ok {
					break
				}

				if msg.context == "" {
					t.Context = nil
				} else {
					t.Context = &msg.context
				}
				t.UpdatedAt = msg.updatedAt
				cmd = m.archivedTaskList.SetItem(msg.listIndex, list.Item(t))
				cmds = append(cmds, cmd)
			}

			if m.activeView == taskDetailsView {
				m.taskDetailsVP.GotoTop()
				m.setContextFSContent(t)
			}
			// to force refresh
			m.contextVPTaskId = 0
		}

	case taskStatusChangedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error deleting task: %s", msg.err)
		} else {
			switch msg.active {
			case true:
				item := m.archivedTaskList.Items()[msg.listIndex]
				oldIndex := m.taskList.Index()

				t, ok := item.(types.Task)
				if !ok {
					break
				}
				t.UpdatedAt = msg.updatedAt
				m.taskList.InsertItem(0, list.Item(t))
				m.taskList.Select(oldIndex + 1)
				m.archivedTaskList.RemoveItem(msg.listIndex)
			case false:
				item := m.taskList.Items()[msg.listIndex]

				t, ok := item.(types.Task)
				if !ok {
					break
				}

				t.UpdatedAt = msg.updatedAt
				m.archivedTaskList.InsertItem(0, list.Item(t))
				m.taskList.RemoveItem(msg.listIndex)
			}
			cmd = m.updateActiveTasksSequence()
			m.updateArchivedTasksIndex()
			cmds = append(cmds, cmd)
		}

	case tasksFetched:
		if msg.err != nil {
			message := "error fetching tasks : " + msg.err.Error()
			m.errorMsg = message
		} else {
			switch msg.active {
			case true:
				taskItems := make([]list.Item, len(msg.tasks))
				for i, t := range msg.tasks {
					taskItems[i] = t
				}
				m.taskList.SetItems(taskItems)
				m.taskList.Select(0)

				tlIndexMap := make(map[uint64]int)
				for i, ti := range m.taskList.Items() {
					t, ok := ti.(types.Task)
					if ok {
						tlIndexMap[t.ID] = i
					}
				}
				m.tlIndexMap = tlIndexMap

			case false:
				archivedTaskItems := make([]list.Item, len(msg.tasks))
				for i, t := range msg.tasks {
					archivedTaskItems[i] = t
				}
				m.archivedTaskList.SetItems(archivedTaskItems)
				m.archivedTaskList.Select(0)
				m.updateArchivedTasksIndex()
			}
		}
	case textEditorClosed:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Something went wrong: %s", msg.err)
			_ = os.Remove(msg.fPath)
			break
		}

		context, err := os.ReadFile(msg.fPath)
		if err != nil {
			break
		}

		err = os.Remove(msg.fPath)
		if err != nil {
			m.errorMsg = fmt.Sprintf("warning: omm failed to remove temporary file: %s", err)
		}

		if len(context) > pers.ContextMaxBytes {
			m.errorMsg = "The content you entered is too large, maybe shorten it"
			// TODO: allow reopening the text editor with the same content again
			break
		}

		if len(context) == 0 && msg.oldContext == nil {
			break
		}

		cmds = append(cmds, updateTaskContext(m.db, msg.taskIndex, msg.taskId, string(context), m.activeTaskList))
	case urlOpenedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error opening url: %s", msg.err)
		}
	case urlsOpenedDarwinMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error opening urls: %s", msg.err)
		}

	case contextWrittenToCBMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Couldn't copy context to clipboard: %s", msg.err)
		} else {
			m.successMsg = "Context copied to clipboard!"
		}
	}

	var viewUpdateCmd tea.Cmd
	switch m.activeView {
	case taskListView:
		m.taskList, viewUpdateCmd = m.taskList.Update(msg)

		if !m.cfg.ShowContext {
			break
		}

		t, ok := m.taskList.SelectedItem().(types.Task)
		if !ok {
			break
		}

		if m.contextVPTaskId == t.ID {
			break
		}

		var detailsToRender string
		switch t.Context {
		case nil:
			detailsToRender = noContextMsg
		default:
			detailsToRender = *t.Context
			switch m.contextMdRenderer {
			case nil:
				break
			default:
				contextGl, err := m.contextMdRenderer.Render(*t.Context)
				if err != nil {
					break
				}
				detailsToRender = contextGl
			}
		}

		m.contextVP.SetContent(detailsToRender)
		m.contextVPTaskId = t.ID

	case archivedTaskListView:
		m.archivedTaskList, viewUpdateCmd = m.archivedTaskList.Update(msg)

		if !m.cfg.ShowContext {
			break
		}

		t, ok := m.archivedTaskList.SelectedItem().(types.Task)
		if !ok {
			break
		}

		if m.contextVPTaskId == t.ID {
			break
		}

		if t.Context != nil {
			if m.contextMdRenderer != nil {
				contextGl, err := m.contextMdRenderer.Render(*t.Context)
				if err != nil {
					m.contextVP.SetContent(*t.Context)
				} else {
					m.contextVP.SetContent(contextGl)
				}
			} else {
				m.contextVP.SetContent(*t.Context)
			}
		} else {
			m.contextVP.SetContent(noContextMsg)
		}
		m.contextVPTaskId = t.ID

	case taskEntryView:
		m.taskInput, viewUpdateCmd = m.taskInput.Update(msg)

	case taskDetailsView:
		m.taskDetailsVP, viewUpdateCmd = m.taskDetailsVP.Update(msg)

	case contextBookmarksView:
		m.taskBMList, viewUpdateCmd = m.taskBMList.Update(msg)

	case prefixSelectionView:
		m.prefixSearchList, viewUpdateCmd = m.prefixSearchList.Update(msg)

	case helpView:
		m.helpVP, viewUpdateCmd = m.helpVP.Update(msg)
	}

	cmds = append(cmds, viewUpdateCmd)

	return m, tea.Batch(cmds...)
}

func (m *model) updateActiveTasksSequence() tea.Cmd {
	sequence := make([]uint64, len(m.taskList.Items()))
	tlIndexMap := make(map[uint64]int)

	for i, ti := range m.taskList.Items() {
		t, ok := ti.(types.Task)
		if ok {
			sequence[i] = t.ID
			tlIndexMap[t.ID] = i
		}
	}

	m.tlIndexMap = tlIndexMap

	return updateTaskSequence(m.db, sequence)
}

func (m *model) updateArchivedTasksIndex() {
	sequence := make([]uint64, len(m.archivedTaskList.Items()))
	tlIndexMap := make(map[uint64]int)

	for i, ti := range m.archivedTaskList.Items() {
		t, ok := ti.(types.Task)
		if ok {
			sequence[i] = t.ID
			tlIndexMap[t.ID] = i
		}
	}

	m.atlIndexMap = tlIndexMap
}

func (m model) isSpaceAvailable() bool {
	return len(m.taskList.Items()) < pers.TaskNumLimit
}

func (m *model) setContextFSContent(task types.Task) {
	var ctx string
	if task.Context != nil {
		ctx = fmt.Sprintf("---\n%s", *task.Context)
	}

	details := fmt.Sprintf(`- summary          :    %s
- created at       :    %s
- last updated at  :    %s

%s
`, task.Summary, task.CreatedAt.Format(timeFormat), task.UpdatedAt.Format(timeFormat), ctx)

	if m.taskDetailsMdRenderer != nil {
		detailsGl, err := m.taskDetailsMdRenderer.Render(details)
		if err == nil {
			m.taskDetailsVP.SetContent(detailsGl)
			return
		}
	}
	m.taskDetailsVP.SetContent(details)
}

func (m model) getTaskUrls() ([]string, bool) {
	var t types.Task
	var ok bool

	switch m.activeView {
	case taskListView:
		t, ok = m.taskList.SelectedItem().(types.Task)
	case archivedTaskListView:
		t, ok = m.archivedTaskList.SelectedItem().(types.Task)
	case taskDetailsView:
		switch m.activeTaskList {
		case activeTasks:
			t, ok = m.taskList.SelectedItem().(types.Task)
		case archivedTasks:
			t, ok = m.archivedTaskList.SelectedItem().(types.Task)
		}
	}
	if !ok {
		return nil, false
	}

	var urls []string
	urls = append(urls, utils.ExtractURLs(m.urlRegex, t.Summary)...)
	if t.Context != nil {
		urls = append(urls, utils.ExtractURLs(m.urlRegex, *t.Context)...)
	}

	return urls, true
}
