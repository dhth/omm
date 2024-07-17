package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/utils"
)

var (
	TaskListDefaultTitle = "omm"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var header string
	var content string
	var context string
	var statusBar string
	var helpMsg string
	var listEmpty bool

	if m.errorMsg != "" && m.successMsg != "" {
		statusBar = fmt.Sprintf("%s%s",
			sBErrMsgStyle.Render(utils.Trim(m.errorMsg, (m.terminalWidth/2)-3)),
			sBSuccessMsgStyle.Render(utils.Trim(m.successMsg, (m.terminalWidth/2)-3)),
		)
	} else if m.errorMsg != "" {
		statusBar = sBErrMsgStyle.Render(m.errorMsg)
	} else if m.successMsg != "" {
		statusBar = sBSuccessMsgStyle.Render(m.successMsg)
	}

	if m.showHelpIndicator && (m.activeView == taskListView || m.activeView == archivedTaskListView) {
		helpMsg = helpMsgStyle.Render("  Press ? for help")
	}

	switch m.activeView {
	case taskListView:
		header = fmt.Sprintf("%s%s", m.tlTitleStyle.Render(m.cfg.TaskListTitle), helpMsg)

		if len(m.taskList.Items()) > 0 {
			content = listStyle.Render(m.taskList.View())
		} else {
			content += fmt.Sprintf("  %s", formStyle.Render("No items. Press a/o to add one.\n"))
			listEmpty = true
		}

		if m.cfg.ListDensity == Compact && len(m.taskList.Items()) <= 9 {
			content += "\n"
		}
	case archivedTaskListView:
		header = fmt.Sprintf("%s%s", m.atlTitleStyle.Render("archived"), helpMsg)
		if len(m.archivedTaskList.Items()) > 0 {
			content = listStyle.Render(m.archivedTaskList.View())
		} else {
			content += fmt.Sprintf("  %s", formStyle.Render("No items. You archive items by pressing ctrl+d.\n"))
			listEmpty = true
		}

		if m.cfg.ListDensity == Compact && len(m.archivedTaskList.Items()) <= 9 {
			content += "\n"
		}
	case taskEntryView:
		if m.taskChange == taskInsert {
			header = taskEntryTitleStyle.Render("enter your task")

			var newTaskPosition string
			if m.taskIndex == 0 {
				newTaskPosition = "at the top"
			} else if m.taskIndex == len(m.taskList.Items()) {
				newTaskPosition = "at the end"
			} else {
				newTaskPosition = fmt.Sprintf("at position %d", m.taskIndex+1)
			}
			content = fmt.Sprintf(`  %s

  %s

  %s

  %s`,
				formHelpStyle.Render(fmt.Sprintf("task will be added %s", newTaskPosition)),
				formHelpStyle.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				formHelpStyle.Render("press <esc> to go back, ⏎ to submit"),
			)

			if m.cfg.ListDensity == Spacious {
				for i := 0; i < m.terminalHeight-13; i++ {
					content += "\n"
				}
			}
		} else if m.taskChange == taskUpdateSummary {
			header = taskEntryTitleStyle.Render("update task")
			content = fmt.Sprintf(`  %s

  %s

  %s`,
				formHelpStyle.Render("omm picks up the prefix in a task summary like 'prefix: do something'\n  and highlights it for you in the task list"),
				m.taskInput.View(),
				formHelpStyle.Render("press <esc> to go back, ⏎ to submit"),
			)
			if m.cfg.ListDensity == Spacious {
				for i := 0; i < m.terminalHeight-11; i++ {
					content += "\n"
				}
			}
		}

	case taskDetailsView:
		var spVal string
		sp := int(m.taskDetailsVP.ScrollPercent() * 100)
		if sp < 100 {
			spVal = helpSectionStyle.Render(fmt.Sprintf("  %d%% ↓", sp))
		}
		header = fmt.Sprintf("%s%s", taskDetailsTitleStyle.Render("task details"), spVal)
		if !m.taskDetailsVPReady {
			context = taskDetailsStyle.Render("Initializing...")
		} else {
			context = taskDetailsStyle.Render(m.taskDetailsVP.View())
		}

		return lipgloss.JoinVertical(lipgloss.Left, headerStyle.Render(header), context, statusBar)

	case contextBookmarksView:
		header = fmt.Sprintf("%s%s", contextBMTitleStyle.Render("Context Bookmarks"), helpMsg)

		content = listStyle.Render(m.contextBMList.View())

	case helpView:
		header = fmt.Sprintf("%s  %s", helpTitleStyle.Render("help"), helpSectionStyle.Render("(scroll with j/k/↓/↑)"))
		if !m.helpVPReady {
			content = helpViewStyle.Render("Initializing...")
		} else {
			content = helpViewStyle.Render(m.helpVP.View())
		}
	}

	var components []string
	components = append(components, headerStyle.Render(header))
	components = append(components, content)

	if !listEmpty && m.cfg.ShowContext && (m.activeView == taskListView || m.activeView == archivedTaskListView) {

		if !m.contextVPReady {
			context = contextStyle.Render("Initializing...")
		} else {
			context = fmt.Sprintf("  %s\n\n%s",
				contextTitleStyle.Render("context"),
				contextStyle.Render(m.contextVP.View()),
			)
		}
		components = append(components, context)
	}

	components = append(components, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, components...)

}
