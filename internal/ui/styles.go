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
	primaryC := lipgloss.Color(thm.Primary)
	secondaryC := lipgloss.Color(thm.Secondary)
	tertiaryC := lipgloss.Color(thm.Tertiary)
	quaternaryC := lipgloss.Color(thm.Quaternary)
	quinaryC := lipgloss.Color(thm.Quinary)
	successC := lipgloss.Color(thm.Success)
	errorC := lipgloss.Color(thm.Error)
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
			Background(quaternaryC),
		helpTitle: titleBase.
			Background(tertiaryC),
		contextTitle: titleBase.
			Background(tertiaryC),
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
			Foreground(errorC),
		statusSuccess: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(successC),
		statusHint: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(mutedC),
		deletePrompt: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(errorC),
		mutedText: mutedStyle,
		activeListTitle: titleBase.
			Background(primaryC),
		archivedListTitle: titleBase.
			Background(secondaryC),
		activeListTitleBar:    listTitleBase.Background(primaryC),
		archivedListTitleBar:  listTitleBase.Background(secondaryC),
		bookmarksListTitleBar: listTitleBase.Background(tertiaryC),
		prefixListTitleBar:    listTitleBase.Background(quinaryC),
		dangerListTitleBar:    listTitleBase.Background(errorC),
	}
}
