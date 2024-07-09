package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/utils"
)

var (
	TaskListDefaultTitle = "omm"
)

func (m model) View() string {
	var content string
	var header string
	var footer string

	if m.message != "" {
		footer = utils.Trim(m.message, 120)
	}

	switch m.activeView {
	case taskListView:
		header = m.tlTitleStyle.Render(m.cfg.TaskListTitle)
		if len(m.taskList.Items()) > 0 {
			content = m.taskList.View()
		} else {
			content = taskInputFormStyle.Render("  No items. Press o to add one.")
		}
	case archivedTaskListView:
		header = m.atlTitleStyle.Render("archived")
		if len(m.archivedTaskList.Items()) > 0 {
			content = m.archivedTaskList.View()
		} else {
			content = taskInputFormStyle.Render("  No items. You archive items by pressing ctrl+d.")
		}
	case taskEntryView:
		header = taskEntryTitleStyle.Render("enter your task")
		content = m.taskInput.View()

	case helpView:
		header = helpTitleStyle.Render("help")
		content = helpSectionStyle.Render(helpMsg)
	}

	if m.showHelpIndicator && (m.activeView == taskListView || m.activeView == archivedTaskListView) {
		header += helpMsgStyle.Render("Press ? for help")
	}

	if m.quitting {
		return ""
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render(header),
		content,
		footerStyle.Render(footer),
	)
}
