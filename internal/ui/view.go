package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var TaskListDefaultTitle = "omm"

func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var content string
	var context string
	var statusBar string
	var listEmpty bool

	if m.showHelpIndicator && (m.activeView != helpView) {
		statusBar += m.styles.statusHint.Render("Press ? for help")
	}

	if m.showDeletePrompt {
		statusBar += m.styles.deletePrompt.Render("press ctrl+x again to delete, any other key to cancel")
	}

	if m.errorMsg != "" && m.successMsg != "" {
		statusBar += fmt.Sprintf("%s%s",
			m.styles.statusError.Render(m.errorMsg),
			m.styles.statusSuccess.Render(m.successMsg),
		)
	} else if m.errorMsg != "" {
		statusBar += m.styles.statusError.Render(m.errorMsg)
	} else if m.successMsg != "" {
		statusBar += m.styles.statusSuccess.Render(m.successMsg)
	}

	switch m.activeView {
	case taskListView:

		if len(m.taskList.Items()) > 0 {
			content = m.styles.listContainer.Render(m.taskList.View())
		} else {
			content = fmt.Sprintf(`
  %s

  %s`, m.styles.activeListTitle.Render(m.cfg.TaskListTitle), m.styles.mutedText.Render("No items. Press a/o to add one.\n"))
			listEmpty = true
		}

	case archivedTaskListView:
		if len(m.archivedTaskList.Items()) > 0 {
			content = m.styles.listContainer.Render(m.archivedTaskList.View())
		} else {
			content = fmt.Sprintf(`
  %s

  %s`, m.styles.archivedListTitle.Render(archivedTitle), m.styles.mutedText.Render("No items. You archive items by pressing ctrl+d.\n"))
			listEmpty = true
		}

	case taskEntryView:
		switch m.taskChange {
		case taskInsert:
			header := m.styles.taskEntryTitle.Render("enter your task")

			var newTaskPosition string
			if m.taskIndex == 0 {
				newTaskPosition = "at the top"
			} else if m.taskIndex == len(m.taskList.Items()) {
				newTaskPosition = "at the end"
			} else {
				newTaskPosition = fmt.Sprintf("at position %d", m.taskIndex+1)
			}
			content = fmt.Sprintf(`
  %s

  %s

  %s

  %s

  %s`,
				header,
				m.styles.mutedText.Render(fmt.Sprintf("task will be added %s", newTaskPosition)),
				m.styles.mutedText.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				m.styles.mutedText.Render("press <esc> to go back, ⏎ to submit"),
			)

			for range m.terminalHeight - 12 {
				content += "\n"
			}
		case taskUpdateSummary:
			header := m.styles.taskEntryTitle.Render("update task")
			content = fmt.Sprintf(`
  %s

  %s

  %s

  %s`,
				header,
				m.styles.mutedText.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				m.styles.mutedText.Render("press <esc> to go back, ⏎ to submit"),
			)
			for range m.terminalHeight - 10 {
				content += "\n"
			}

		}

	case taskDetailsView:
		var spVal string
		sp := int(m.taskDetailsVP.ScrollPercent() * 100)
		if sp < 100 {
			spVal = m.styles.statusHint.Render(fmt.Sprintf("  %d%% ↓", sp))
		}
		header := fmt.Sprintf("%s%s", m.styles.taskDetailsTitle.Render("task details"), spVal)
		if !m.taskDetailsVPReady {
			content = m.styles.sectionHeader.Render(header) + "\n" + "Initializing..."
		} else {
			content = m.styles.sectionHeader.Render(header) + "\n" + m.taskDetailsVP.View()
		}

	case contextBookmarksView:
		content = m.styles.listContainer.Render(m.taskBMList.View())

	case prefixSelectionView:
		content = m.styles.listContainer.Render(m.prefixSearchList.View())

	case helpView:
		header := fmt.Sprintf(`
  %s  %s

`, m.styles.helpTitle.Render("help"), m.styles.statusHint.Render("(scroll with j/k/↓/↑)"))
		if !m.helpVPReady {
			content = "Initializing..."
		} else {
			content = header + m.helpVP.View()
		}
	}

	var components []string
	components = append(components, content)

	if !listEmpty && m.cfg.ShowContext && (m.activeView == taskListView || m.activeView == archivedTaskListView) {

		if !m.contextVPReady {
			context = "Initializing..."
		} else {
			context = fmt.Sprintf("  %s\n\n%s",
				m.styles.contextTitle.Render("context"),
				m.contextVP.View(),
			)
		}
		components = append(components, context)
	}

	components = append(components, statusBar)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, components...))
	v.AltScreen = true

	return v
}
