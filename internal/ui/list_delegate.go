package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newListDelegate(color lipgloss.Color, showDesc bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = showDesc

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	d.Styles.FilterMatch = lipgloss.NewStyle()

	return d
}
