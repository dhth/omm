package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	TaskListColor          = "#fe8019"
	ArchivedTLColor        = "#fabd2f"
	taskEntryTitleColor    = "#b8bb26"
	taskListHeaderColor    = "#928374"
	taskInputFormColor     = "#928374"
	helpMsgColor           = "#928374"
	helpViewTitleColor     = "#83a598"
	helpTitleColor         = "#83a598"
	helpSectionColor       = "#928374"
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

	taskEntryTitleStyle = titleStyle.
				Background(lipgloss.Color(taskEntryTitleColor))

	helpTitleStyle = titleStyle.
			Background(lipgloss.Color(helpTitleColor))

	headerStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2)

	footerStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color(footerColor))

	taskInputFormStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(taskInputFormColor))

	helpSectionStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(helpSectionColor))

	helpMsgStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color(helpMsgColor))
)
