package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	TaskListColor          = "#fe8019"
	ArchivedTLColor        = "#fabd2f"
	contextTitleColor      = "#8ec07c"
	taskEntryTitleColor    = "#b8bb26"
	taskDetailsTitleColor  = "#d3869b"
	taskListHeaderColor    = "#928374"
	contextTextColor       = "#928374"
	formHelpColor          = "#928374"
	formColor              = "#928374"
	helpMsgColor           = "#928374"
	helpViewTitleColor     = "#83a598"
	helpTitleColor         = "#83a598"
	helpSectionColor       = "#928374"
	statusBarColor         = "#fb4934"
	footerColor            = "#928374"
)

var (
	itemStyle = lipgloss.NewStyle().PaddingLeft(2)

	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true).
			Background(lipgloss.Color(TaskListColor)).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	listStyle = lipgloss.NewStyle().PaddingBottom(1)

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

	statusBarStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color(statusBarColor))

	contextTextStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(contextTextColor))

	formStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formColor))

	formHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formHelpColor))

	helpSectionStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(helpSectionColor))

	helpMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpMsgColor))
)
