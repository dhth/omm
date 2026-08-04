package ui

import (
	"charm.land/lipgloss/v2"
	"github.com/dhth/omm/internal/ui/theme"
)

type styles struct {
	listContainer         lipgloss.Style
	taskEntryTitle        lipgloss.Style
	helpTitle             lipgloss.Style
	contextTitle          lipgloss.Style
	taskDetailsTitle      lipgloss.Style
	sectionHeader         lipgloss.Style
	statusBar             lipgloss.Style
	statusError           lipgloss.Style
	statusSuccess         lipgloss.Style
	statusHint            lipgloss.Style
	deletePrompt          lipgloss.Style
	mutedText             lipgloss.Style
	activeListTitle       lipgloss.Style
	archivedListTitle     lipgloss.Style
	activeListTitleBar    lipgloss.Style
	archivedListTitleBar  lipgloss.Style
	bookmarksListTitleBar lipgloss.Style
	prefixListTitleBar    lipgloss.Style
	dangerListTitleBar    lipgloss.Style
}

func newStyles(thm theme.Theme) styles {
	bg := lipgloss.Color(thm.Background)
	accent1C := lipgloss.Color(thm.Accent1)
	accent2C := lipgloss.Color(thm.Accent2)
	accent3C := lipgloss.Color(thm.Accent3)
	accent4C := lipgloss.Color(thm.Accent4)
	accent5C := lipgloss.Color(thm.Accent5)
	successC := lipgloss.Color(thm.Success)
	dangerC := lipgloss.Color(thm.Danger)
	mutedC := lipgloss.Color(thm.Muted)

	mutedStyle := lipgloss.NewStyle().Foreground(mutedC)

	titleBase := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Bold(true).
		Foreground(bg)

	listTitleBase := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(bg).
		Bold(true)

	return styles{
		listContainer: lipgloss.NewStyle().PaddingBottom(1).PaddingTop(1),
		taskEntryTitle: titleBase.
			Background(accent4C),
		helpTitle: titleBase.
			Background(accent3C),
		contextTitle: titleBase.
			Background(accent3C),
		taskDetailsTitle: titleBase.
			Background(successC),
		sectionHeader: lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2),
		statusBar: lipgloss.NewStyle().
			PaddingLeft(2),
		statusError: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(dangerC),
		statusSuccess: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(successC),
		statusHint: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(mutedC),
		deletePrompt: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(dangerC),
		mutedText: mutedStyle,
		activeListTitle: titleBase.
			Background(accent1C),
		archivedListTitle: titleBase.
			Background(accent2C),
		activeListTitleBar:    listTitleBase.Background(accent1C),
		archivedListTitleBar:  listTitleBase.Background(accent2C),
		bookmarksListTitleBar: listTitleBase.Background(accent3C),
		prefixListTitleBar:    listTitleBase.Background(accent5C),
		dangerListTitleBar:    listTitleBase.Background(dangerC),
	}
}
