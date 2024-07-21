package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	TaskListColor          = "#fe8019"
	ArchivedTLColor        = "#fabd2f"
	contextBMColor         = "#83a598"
	contextTitleColor      = "#8ec07c"
	taskEntryTitleColor    = "#b8bb26"
	taskDetailsTitleColor  = "#d3869b"
	taskListHeaderColor    = "#928374"
	taskDetailsColor       = "#bdae93"
	contextColor           = "#928374"
	formHelpColor          = "#928374"
	formColor              = "#928374"
	helpMsgColor           = "#83a598"
	helpViewTitleColor     = "#83a598"
	helpTitleColor         = "#83a598"
	helpHeaderColor        = "#83a598"
	helpSectionColor       = "#bdae93"
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

	taskDetailsStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(taskDetailsColor))

	contextStyle = taskDetailsStyle.
			Foreground(lipgloss.Color(contextColor))

	formStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formColor))

	formHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formHelpColor))

	helpViewStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	helpStyle = lipgloss.NewStyle()

	helpHeadingStyle = helpStyle.
				Bold(true).
				Foreground(lipgloss.Color(helpHeaderColor))

	helpSectionStyle = helpStyle.
				Foreground(lipgloss.Color(helpSectionColor))

	helpSubHeadingStyle = helpSectionStyle.
				Bold(true)

	helpMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpMsgColor))
)
