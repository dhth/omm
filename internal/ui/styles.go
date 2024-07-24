package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	TaskListColor          = "#fe8019"
	ArchivedTLColor        = "#fabd2f"
	contextBMColor         = "#83a598"
	prefixSearchColor      = "#d3896b"
	contextTitleColor      = "#8ec07c"
	taskEntryTitleColor    = "#b8bb26"
	taskDetailsTitleColor  = "#d3869b"
	taskListHeaderColor    = "#928374"
	formHelpColor          = "#928374"
	formColor              = "#928374"
	helpMsgColor           = "#928374"
	helpViewTitleColor     = "#83a598"
	helpTitleColor         = "#83a598"
	sBSuccessMsgColor      = "#d3869b"
	sBErrMsgColor          = "#fb4934"
	footerColor            = "#928374"
)

var (
	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true).
			Background(lipgloss.Color(TaskListColor)).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	listStyle = lipgloss.NewStyle().PaddingBottom(1).PaddingTop(1)

	taskEntryTitleStyle = titleStyle.
				Background(lipgloss.Color(taskEntryTitleColor))

	helpTitleStyle = titleStyle.
			Background(lipgloss.Color(helpTitleColor))

	contextTitleStyle = titleStyle.
				Background(lipgloss.Color(contextTitleColor))

	taskDetailsTitleStyle = titleStyle.
				Background(lipgloss.Color(taskDetailsTitleColor))

	headerStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2)

	statusBarMsgStyle = lipgloss.NewStyle().
				PaddingLeft(2)

	sBErrMsgStyle = statusBarMsgStyle.
			Foreground(lipgloss.Color(sBErrMsgColor))

	sBSuccessMsgStyle = statusBarMsgStyle.
				Foreground(lipgloss.Color(sBSuccessMsgColor))

	formStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formColor))

	formHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formHelpColor))

	helpMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpMsgColor))
)
