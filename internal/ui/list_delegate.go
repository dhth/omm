package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newListDelegate(color lipgloss.Color, showDesc bool, spacing int) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = showDesc
	d.SetSpacing(spacing)

	d.Styles.NormalTitle = d.Styles.
		NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ffffff"})

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	d.Styles.FilterMatch = lipgloss.NewStyle()

	return d
}
