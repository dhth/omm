package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/utils"
)

var (
	TaskListDefaultTitle      = "omm"
	taskDetailsWordWrap       = 120
	contextWordWrapUpperLimit = 160
)

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var content string
	var context string
	var statusBar string
	var listEmpty bool

	if m.showHelpIndicator && (m.activeView != helpView) {
		statusBar += helpMsgStyle.Render("Press ? for help")
	}

	if m.showDeletePrompt {
		statusBar += promptStyle.Render("press ctrl+x again to delete, any other key to cancel")
	}

	if m.errorMsg != "" && m.successMsg != "" {
		statusBar += fmt.Sprintf("%s%s",
			sBErrMsgStyle.Render(utils.Trim(m.errorMsg, (m.terminalWidth/2)-3)),
			sBSuccessMsgStyle.Render(utils.Trim(m.successMsg, (m.terminalWidth/2)-3)),
		)
	} else if m.errorMsg != "" {
		statusBar += sBErrMsgStyle.Render(m.errorMsg)
	} else if m.successMsg != "" {
		statusBar += sBSuccessMsgStyle.Render(m.successMsg)
	}

	switch m.activeView {
	case taskListView:

		if len(m.taskList.Items()) > 0 {
			content = listStyle.Render(m.taskList.View())
		} else {
			content = fmt.Sprintf(`
  %s

  %s`, m.tlTitleStyle.Render(m.cfg.TaskListTitle), formStyle.Render("No items. Press a/o to add one.\n"))
			listEmpty = true
		}

	case archivedTaskListView:
		if len(m.archivedTaskList.Items()) > 0 {
			content = listStyle.Render(m.archivedTaskList.View())
		} else {
			content = fmt.Sprintf(`
  %s

  %s`, m.atlTitleStyle.Render(archivedTitle), formStyle.Render("No items. You archive items by pressing ctrl+d.\n"))
			listEmpty = true
		}

	case taskEntryView:
		if m.taskChange == taskInsert {
			header := taskEntryTitleStyle.Render("enter your task")

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
				formHelpStyle.Render(fmt.Sprintf("task will be added %s", newTaskPosition)),
				formHelpStyle.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				formHelpStyle.Render("press <esc> to go back, ⏎ to submit"),
			)

			for i := 0; i < m.terminalHeight-12; i++ {
				content += "\n"
			}
		} else if m.taskChange == taskUpdateSummary {
			header := taskEntryTitleStyle.Render("update task")
			content = fmt.Sprintf(`
  %s

  %s

  %s

  %s`,
				header,
				formHelpStyle.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				formHelpStyle.Render("press <esc> to go back, ⏎ to submit"),
			)
			for i := 0; i < m.terminalHeight-10; i++ {
				content += "\n"
			}
		}

	case taskDetailsView:
		var spVal string
		sp := int(m.taskDetailsVP.ScrollPercent() * 100)
		if sp < 100 {
			spVal = helpMsgStyle.Render(fmt.Sprintf("  %d%% ↓", sp))
		}
		header := fmt.Sprintf("%s%s", taskDetailsTitleStyle.Render("task details"), spVal)
		if !m.taskDetailsVPReady {
			content = headerStyle.Render(header) + "\n" + "Initializing..."
		} else {
			content = headerStyle.Render(header) + "\n" + m.taskDetailsVP.View()
		}

	case contextBookmarksView:
		content = listStyle.Render(m.taskBMList.View())

	case prefixSelectionView:
		content = listStyle.Render(m.prefixSearchList.View())

	case helpView:
		header := fmt.Sprintf(`
  %s  %s

`, helpTitleStyle.Render("help"), helpMsgStyle.Render("(scroll with j/k/↓/↑)"))
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
				contextTitleStyle.Render("context"),
				m.contextVP.View(),
			)
		}
		components = append(components, context)
	}

	components = append(components, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, components...)
}
