package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
)

const (
	noSpaceAvailableMsg = "Task list is at capacity. Archive/delete tasks using ctrl+d/ctrl+x."
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""

	if m.activeView == taskEntryView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "esc", "ctrl+c":
				m.activeView = taskListView
			case "enter":
				taskSummary := m.taskInput.Value()
				taskSummary = strings.TrimSpace(taskSummary)

				if taskSummary == "" {
					m.activeView = taskListView
					break
				}

				switch m.taskChange {
				case taskInsert:
					now := time.Now()
					cmd = createTask(m.db, taskSummary, now, now)
					cmds = append(cmds, cmd)
					m.taskInput.Reset()
					m.activeView = taskListView
				case taskUpdateSummary:
					cmd = updateTaskSummary(m.db, m.taskIndex, m.taskId, taskSummary)
					cmds = append(cmds, cmd)
					m.taskInput.Reset()
					m.activeView = taskListView
				}
			}
		}

		m.taskInput, cmd = m.taskInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.taskList.SetWidth(msg.Width - 2)
		m.archivedTaskList.SetWidth(msg.Width - 2)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "esc", "q", "ctrl+c":
			if m.activeView != taskListView {
				m.activeView = taskListView
				break
			}
			m.quitting = true
			return m, tea.Quit

		case "?":
			m.activeView = helpView

		case "tab", "shift+tab":
			switch m.activeView {
			case taskListView:
				m.activeView = archivedTaskListView
			case archivedTaskListView:
				m.activeView = taskListView
			}

		case "2", "3", "4", "5", "6", "7", "8", "9":
			if m.activeView != taskListView {
				break
			}

			keyNum, err := strconv.Atoi(keypress)
			if err != nil {
				m.message = "Something went horribly wrong"
				break
			}

			if m.taskList.Index() == 0 && keyNum == 1 {
				break
			}

			index := (m.taskList.Paginator.Page * m.taskList.Paginator.PerPage) + (keyNum - 1)

			if index >= len(m.taskList.Items()) {
				m.message = "There is no item for this index"
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
			if m.activeView == taskEntryView {
				break
			}

			if !m.isSpaceAvailable() {
				m.message = noSpaceAvailableMsg
				break
			}

			m.taskIndex = 0
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "O":
			if m.activeView == taskEntryView {
				break
			}

			if !m.isSpaceAvailable() {
				m.message = noSpaceAvailableMsg
				break
			}

			m.taskIndex = m.taskList.Index()
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "a", "o":
			if m.activeView == taskEntryView {
				break
			}

			if !m.isSpaceAvailable() {
				m.message = noSpaceAvailableMsg
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
			if m.activeView == taskEntryView {
				break
			}

			if !m.isSpaceAvailable() {
				m.message = noSpaceAvailableMsg
				break
			}

			m.taskIndex = len(m.taskList.Items())
			m.taskInput.Reset()
			m.taskInput.Focus()
			m.taskChange = taskInsert
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "J":
			if m.activeView == taskEntryView {
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
			if m.activeView == taskEntryView {
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
			task, ok := listItem.(types.Task)
			if !ok {
				m.message = "Something went wrong"
				break
			}

			m.taskInput.SetValue(task.Summary)
			m.taskInput.Focus()
			m.taskIndex = index
			m.taskId = task.ID
			m.taskChange = taskUpdateSummary
			m.activeView = taskEntryView
			return m, tea.Batch(cmds...)

		case "ctrl+d":
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				listItem := m.taskList.SelectedItem()
				index := m.taskList.Index()
				task, ok := listItem.(types.Task)
				if !ok {
					m.message = "Something went wrong"
					break
				}

				cmd = changeTaskStatus(m.db, index, task.ID, false, time.Now())
				cmds = append(cmds, cmd)

			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					break
				}

				listItem := m.archivedTaskList.SelectedItem()
				index := m.archivedTaskList.Index()
				task, ok := listItem.(types.Task)
				if !ok {
					m.message = "Something went wrong"
					break
				}

				cmd = changeTaskStatus(m.db, index, task.ID, true, time.Now())
				cmds = append(cmds, cmd)
			}

		case "ctrl+x":
			switch m.activeView {
			case taskListView:
				if len(m.taskList.Items()) == 0 {
					break
				}

				index := m.taskList.Index()
				task, ok := m.taskList.SelectedItem().(types.Task)
				if !ok {
					m.message = "Something went wrong"
					break
				}
				cmd = deleteTask(m.db, task.ID, index, true)
				cmds = append(cmds, cmd)
			case archivedTaskListView:
				if len(m.archivedTaskList.Items()) == 0 {
					break
				}

				index := m.archivedTaskList.Index()
				task, ok := m.archivedTaskList.SelectedItem().(types.Task)
				if !ok {
					m.message = "Something went wrong"
					break
				}

				cmd = deleteTask(m.db, task.ID, index, false)
				cmds = append(cmds, cmd)
			}

		case "enter":
			if m.activeView == taskListView {
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
			}
		}

	case HideHelpMsg:
		m.showHelpIndicator = false

	case taskCreatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error creating task: %s", msg.err)
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
		m.taskList.Select(m.taskIndex)
		cmds = append(cmds, cmd)

		cmd = m.updateTaskSequence()
		cmds = append(cmds, cmd)

	case taskDeletedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error deleting task: %s", msg.err)
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
			m.message = fmt.Sprintf("Error updating task sequence: %s", msg.err)
		}

	case taskSummaryUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task: %s", msg.err)
		} else {
			listItem := m.taskList.Items()[msg.listIndex]
			task, ok := listItem.(types.Task)
			if !ok {
				break
			}

			task.Summary = msg.taskSummary
			cmd = m.taskList.SetItem(msg.listIndex, list.Item(task))
			cmds = append(cmds, cmd)
		}

	case taskStatusChangedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error deleting task: %s", msg.err)
		} else {
			switch msg.active {
			case true:
				item := m.archivedTaskList.Items()[msg.listIndex]
				m.taskList.InsertItem(0, item)
				m.archivedTaskList.RemoveItem(msg.listIndex)
			case false:
				item := m.taskList.Items()[msg.listIndex]
				m.archivedTaskList.InsertItem(0, item)
				m.taskList.RemoveItem(msg.listIndex)
			}
			cmd = m.updateTaskSequence()
			cmds = append(cmds, cmd)
		}

	case tasksFetched:
		if msg.err != nil {
			message := "error fetching tasks : " + msg.err.Error()
			m.message = message
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
	}

	if len(m.taskList.Items()) > 9 {
		m.taskList.SetHeight(11)
	}

	switch m.activeView {
	case taskListView:
		m.taskList, cmd = m.taskList.Update(msg)
	case archivedTaskListView:
		m.archivedTaskList, cmd = m.archivedTaskList.Update(msg)
	case taskEntryView:
		m.taskInput, cmd = m.taskInput.Update(msg)
	}
	cmds = append(cmds, cmd)

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
