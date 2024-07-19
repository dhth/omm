package ui

import (
	"fmt"
	"os"
	"strconv"
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
	noContextMsg          = "∅"
	viewPortMoveLineCount = 3
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.successMsg = ""
	m.errorMsg = ""

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
			}
		}

		m.taskInput, cmd = m.taskInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := listStyle.GetFrameSize()
		_, h2 := headerStyle.GetFrameSize()
		_, h3 := statusBarMsgStyle.GetFrameSize()
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		m.taskList.SetWidth(msg.Width - w)
		m.archivedTaskList.SetWidth(msg.Width - 2)
		m.contextBMList.SetWidth(msg.Width - 2)
		m.contextBMList.SetHeight(msg.Height - h - 4)

		var listHeight int
		if m.cfg.ShowContext {
			listHeight = msg.Height/2 - h
		} else {
			listHeight = msg.Height - h - 4
		}

		if m.cfg.ListDensity == Spacious {
			m.taskList.SetHeight(listHeight)
			m.archivedTaskList.SetHeight(listHeight)
		}

		var contextHeight int
		if m.cfg.ListDensity == Compact {
			contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 6
		} else {
			contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 5
		}

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

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width-3, m.terminalHeight-4)
			m.helpVP.SetContent(helpStr)
			m.helpVP.KeyMap.Up.SetEnabled(false)
			m.helpVP.KeyMap.Down.SetEnabled(false)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width - 3
			m.helpVP.Height = m.terminalHeight - 4
		}

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "Q":
			m.quitting = true
			if m.cfg.Guide {
				_ = os.Remove(m.cfg.DBPath)
			}
			return m, tea.Quit

		case "esc", "q", "ctrl+c":
			av := m.activeView

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

			m.quitting = true
			if m.cfg.Guide {
				_ = os.Remove(m.cfg.DBPath)
			}
			return m, tea.Quit

		case "?":
			if m.activeView == taskDetailsView {
				break
			}

			if m.activeView == helpView {
				m.activeView = m.lastActiveView
				break
			}
			m.lastActiveView = m.activeView
			m.activeView = helpView

		case "tab", "shift+tab":
			m.lastActiveView = m.activeView
			switch m.activeView {
			case taskListView:
				m.activeView = archivedTaskListView
				m.activeTaskList = archivedTasks
			case archivedTaskListView:
				m.activeView = taskListView
				m.activeTaskList = activeTasks
			}

		case "2", "3", "4", "5", "6", "7", "8", "9":
			if m.cfg.ListDensity != Compact {
				break
			}

			if m.activeView != taskListView {
				break
			}

			keyNum, err := strconv.Atoi(keypress)
			if err != nil {
				m.errorMsg = "Something went horribly wrong"
				break
			}

			if m.taskList.Index() == 0 && keyNum == 1 {
				break
			}

			index := (m.taskList.Paginator.Page * m.taskList.Paginator.PerPage) + (keyNum - 1)

			if index >= len(m.taskList.Items()) {
				m.errorMsg = "There is no item for this index"
				break
			}
			listItem := m.taskList.Items()[index]

			m.taskList.RemoveItem(index)
			cmd = m.taskList.InsertItem(0, listItem)
			cmds = append(cmds, cmd)
			m.taskList.Select(0)

			cmd = m.updateTaskSequence()
			cmds = append(cmds, cmd)

		case "I":
			if m.activeView != taskListView {
				break
			}

			if !m.isSpaceAvailable() {
				m.errorMsg = noSpaceAvailableMsg
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

			m.taskIndex = len(m.taskList.Items())
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "j":
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

		case "k":
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

			cmd = m.updateTaskSequence()
			cmds = append(cmds, cmd)

		case "K":
			if m.activeView != taskListView {
				break
			}

			if m.taskList.IsFiltered() {
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

			cmd = m.updateTaskSequence()
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

			cmds = append(cmds, fetchTasks(m.db, true, pers.TaskNumLimit))
			cmds = append(cmds, fetchTasks(m.db, false, pers.TaskNumLimit))

		case "ctrl+d":
			switch m.activeView {
			case taskListView:
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

				cmd = changeTaskStatus(m.db, index, t.ID, false, time.Now())
				cmds = append(cmds, cmd)

			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
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
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				index := m.taskList.Index()
				t, ok := m.taskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}
				cmd = deleteTask(m.db, t.ID, index, true)
				cmds = append(cmds, cmd)

			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					break
				}

				index := m.archivedTaskList.Index()
				task, ok := m.archivedTaskList.SelectedItem().(types.Task)
				if !ok {
					m.errorMsg = "Something went wrong"
					break
				}

				cmd = deleteTask(m.db, task.ID, index, false)
				cmds = append(cmds, cmd)
			}

		case "enter":
			if m.activeView != taskListView && m.activeView != contextBookmarksView {
				break
			}
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				if m.taskList.Index() == 0 {
					break
				}

				index := m.taskList.Index()
				listItem := m.taskList.SelectedItem()
				m.taskList.RemoveItem(index)
				cmd = m.taskList.InsertItem(0, listItem)
				cmds = append(cmds, cmd)
				m.taskList.Select(0)

				cmd = m.updateTaskSequence()
				cmds = append(cmds, cmd)
			case contextBookmarksView:
				url := m.contextBMList.SelectedItem().FilterValue()
				cmds = append(cmds, openURL(url))
			}

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

			w, h := listStyle.GetFrameSize()
			var taskList list.Model
			var archivedTaskList list.Model
			_, h2 := headerStyle.GetFrameSize()
			_, h3 := statusBarMsgStyle.GetFrameSize()

			tlIndex := m.taskList.Index()
			atlIndex := m.archivedTaskList.Index()

			var listHeight int
			if m.cfg.ShowContext {
				listHeight = m.terminalHeight/2 - h
			} else {
				listHeight = m.terminalHeight - h - 4
			}

			switch m.cfg.ListDensity {
			case Compact:
				taskList = list.New(m.taskList.Items(),
					newTaskListDelegate(lipgloss.Color(m.cfg.TaskListColor)),
					m.terminalWidth-w,
					listHeight,
				)
				taskList.SetShowStatusBar(true)

				archivedTaskList = list.New(m.archivedTaskList.Items(),
					newTaskListDelegate(lipgloss.Color(m.cfg.ArchivedTaskListColor)),
					m.terminalWidth-w,
					listHeight,
				)
				archivedTaskList.SetShowStatusBar(true)

				m.cfg.ListDensity = Spacious

			case Spacious:
				taskList = list.New(m.taskList.Items(),
					itemDelegate{selStyle: m.tlSelStyle},
					m.terminalWidth-w,
					compactListHeight,
				)
				taskList.SetShowStatusBar(false)

				archivedTaskList = list.New(m.archivedTaskList.Items(),
					itemDelegate{selStyle: m.tlSelStyle},
					m.terminalWidth-w,
					compactListHeight,
				)
				m.cfg.ListDensity = Compact
				archivedTaskList.SetShowStatusBar(false)
			}

			taskList.SetShowTitle(false)
			taskList.SetFilteringEnabled(false)
			taskList.SetShowHelp(false)
			taskList.DisableQuitKeybindings()
			taskList.Styles.Title = m.taskList.Styles.Title
			taskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
			taskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

			m.taskList = taskList
			m.taskList.Select(tlIndex)

			archivedTaskList.SetShowTitle(false)
			archivedTaskList.SetFilteringEnabled(false)
			archivedTaskList.SetShowHelp(false)
			archivedTaskList.DisableQuitKeybindings()
			archivedTaskList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
			archivedTaskList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

			m.archivedTaskList = archivedTaskList
			m.archivedTaskList.Select(atlIndex)

			var contextHeight int
			if m.cfg.ListDensity == Compact {
				contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 6
			} else {
				contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 5
			}
			m.contextVP.Height = contextHeight

		case "C":
			if m.activeView != taskListView && m.activeView != archivedTaskListView {
				break
			}

			m.cfg.ShowContext = !m.cfg.ShowContext

			_, h := listStyle.GetFrameSize()
			var listHeight int
			if m.cfg.ListDensity == Spacious {
				switch m.cfg.ShowContext {
				case true:
					listHeight = m.terminalHeight/2 - h
				case false:
					listHeight = m.terminalHeight - h - 4
				}
				m.taskList.SetHeight(listHeight)
				m.archivedTaskList.SetHeight(listHeight)
			}

			_, h2 := headerStyle.GetFrameSize()
			_, h3 := statusBarMsgStyle.GetFrameSize()

			var contextHeight int
			if m.cfg.ListDensity == Compact {
				contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 6
			} else {
				contextHeight = m.terminalHeight - m.taskList.Height() - h2 - h3 - 5
			}
			m.contextVP.Height = contextHeight

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
			m.contextBMList.SetItems(bmItems)
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

		entry := list.Item(types.Task{
			ID:        msg.id,
			Summary:   msg.taskSummary,
			Active:    true,
			CreatedAt: msg.createdAt,
			UpdatedAt: msg.updatedAt,
		})
		cmd = m.taskList.InsertItem(m.taskIndex, entry)
		cmds = append(cmds, cmd)
		m.taskList.Select(m.taskIndex)

		cmd = m.updateTaskSequence()
		cmds = append(cmds, cmd)

	case taskDeletedMsg:
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error deleting task: %s", msg.err)
			break
		}

		switch msg.active {
		case true:
			m.taskList.RemoveItem(msg.listIndex)
			cmd = m.updateTaskSequence()
			cmds = append(cmds, cmd)
		case false:
			m.archivedTaskList.RemoveItem(msg.listIndex)
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
			cmd = m.updateTaskSequence()
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
			case false:
				archivedTaskItems := make([]list.Item, len(msg.tasks))
				for i, t := range msg.tasks {
					archivedTaskItems[i] = t
				}
				m.archivedTaskList.SetItems(archivedTaskItems)
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

	if m.cfg.ListDensity == Compact {
		if len(m.taskList.Items()) > 9 {
			m.taskList.SetHeight(compactListHeight + 1)
		} else {
			m.taskList.SetHeight(compactListHeight)
		}
		if len(m.archivedTaskList.Items()) > 9 {
			m.archivedTaskList.SetHeight(compactListHeight + 1)
		} else {
			m.archivedTaskList.SetHeight(compactListHeight)
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

		if t.Context != nil {
			m.contextVP.SetContent(*t.Context)
		} else {
			m.contextVP.SetContent(noContextMsg)
		}
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
			m.contextVP.SetContent(*t.Context)
		} else {
			m.contextVP.SetContent(noContextMsg)
		}
		m.contextVPTaskId = t.ID

	case taskEntryView:
		m.taskInput, viewUpdateCmd = m.taskInput.Update(msg)

	case taskDetailsView:
		m.taskDetailsVP, viewUpdateCmd = m.taskDetailsVP.Update(msg)

	case contextBookmarksView:
		m.contextBMList, viewUpdateCmd = m.contextBMList.Update(msg)

	case helpView:
		m.helpVP, viewUpdateCmd = m.helpVP.Update(msg)
	}

	cmds = append(cmds, viewUpdateCmd)

	return m, tea.Batch(cmds...)
}

func (m model) updateTaskSequence() tea.Cmd {
	sequence := make([]uint64, len(m.taskList.Items()))
	for i, ti := range m.taskList.Items() {
		t, ok := ti.(types.Task)
		if ok {
			sequence[i] = t.ID
		}
	}

	return updateTaskSequence(m.db, sequence)
}

func (m model) isSpaceAvailable() bool {
	return len(m.taskList.Items()) < pers.TaskNumLimit
}

func (m *model) setContextFSContent(task types.Task) {
	var ctx string
	if task.Context != nil {
		ctx = fmt.Sprintf("\n===\n\n%s", *task.Context)
	}

	details := fmt.Sprintf(`summary               %s
created at            %s
last updated at       %s
%s
`, task.Summary, task.CreatedAt.Format(timeFormat), task.UpdatedAt.Format(timeFormat), ctx)

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
