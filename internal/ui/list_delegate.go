package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newTaskListDelegate(color lipgloss.Color) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}

func newContextURLListDel(color lipgloss.Color) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	d.ShowDescription = false
	d.SetHeight(1)

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
