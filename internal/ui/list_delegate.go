package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/ui/theme"
	"github.com/dhth/omm/internal/utils"
	"github.com/dustin/go-humanize"
)

const (
	spaciousPrefixPadding = 80
	createdAtPadding      = 40
)

type compactItemDelegate struct {
	selStyle     lipgloss.Style
	prefixColors []string
}

type spaciousTaskItemDelegate struct {
	selStyle           lipgloss.Style
	secondaryTextStyle lipgloss.Style
	prefixColors       []string
}

func (d compactItemDelegate) Height() int { return 1 }

func (d compactItemDelegate) Spacing() int { return 1 }

func (d compactItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d spaciousTaskItemDelegate) Height() int { return 2 }

func (d spaciousTaskItemDelegate) Spacing() int { return 1 }

func (d spaciousTaskItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d compactItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(types.Task)
	if !ok {
		return
	}

	prefix, sc, hp := t.GetPrefixAndSummaryContent()

	if hp {
		prefix = lipgloss.NewStyle().
			Foreground(getColorForString(prefix, d.prefixColors)).
			Bold(true).
			Render(utils.RightPadTrim(prefix, prefixPadding, true))
	}
	var hasContext string
	if t.Context != nil {
		hasContext = "(c)"
	}

	sr := d.selStyle.Render
	var str string
	if index == m.Index() {
		str = fmt.Sprintf("%s%s%s%s", sr("▎ "), prefix, sr(utils.RightPadTrim(sc, taskSummaryWidth-prefixPadding, true)), sr(hasContext))
	} else {
		str = fmt.Sprintf("%s%s%s%s", "  ", prefix, utils.RightPadTrim(sc, taskSummaryWidth-prefixPadding, true), hasContext)
	}

	fmt.Fprint(w, str)
}

func (d spaciousTaskItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(types.Task)
	if !ok {
		return
	}

	prefix, sc, hasPrefix := t.GetPrefixAndSummaryContent()
	if hasPrefix {
		prefix = lipgloss.NewStyle().
			Foreground(getColorForString(prefix, d.prefixColors)).
			Bold(true).
			Render(utils.RightPadTrim(prefix, spaciousPrefixPadding, true))
	} else {
		prefix = strings.Repeat(" ", spaciousPrefixPadding)
	}

	createdAtTs := humanize.Time(t.CreatedAt)
	if time.Since(t.CreatedAt).Seconds() < 60 {
		createdAtTs = "just now"
	}

	createdAt := d.secondaryTextStyle.Render(utils.RightPadTrim(fmt.Sprintf("created %s", createdAtTs), createdAtPadding, true))

	hasContext := ""
	if t.Context != nil {
		hasContext = d.secondaryTextStyle.Render("(c)")
	}

	desc := fmt.Sprintf("%s%s%s", prefix, createdAt, hasContext)
	title := utils.RightPadTrim(sc, taskSummaryWidth-2, true)
	desc = utils.RightPadTrim(desc, taskSummaryWidth-2, true)

	sr := d.selStyle.Render
	if index == m.Index() {
		fmt.Fprintf(w, "%s%s\n%s%s", sr("▎ "), sr(title), sr("▎ "), sr(desc))
		return
	}

	fmt.Fprintf(w, "  %s\n  %s", title, desc)
}

func newTaskListDelegate(thm theme.Theme, density ListDensityType, listType taskListType) list.ItemDelegate {
	selectionColor := lipgloss.Color(thm.Primary)
	if listType == archivedTasks {
		selectionColor = lipgloss.Color(thm.Secondary)
	}

	selectionStyle := lipgloss.NewStyle().Foreground(selectionColor)

	switch density {
	case Spacious:
		secondaryTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Muted))
		return spaciousTaskItemDelegate{selectionStyle, secondaryTextStyle, thm.PrefixColors}
	default:
		return compactItemDelegate{selectionStyle, thm.PrefixColors}
	}
}

func newBookmarksListDelegate(thm theme.Theme) list.ItemDelegate {
	return newSpaciousListDelegate(lipgloss.Color(thm.Tertiary), false, 1)
}

func newPrefixSearchListDelegate(thm theme.Theme) list.ItemDelegate {
	return newSpaciousListDelegate(lipgloss.Color(thm.Quinary), false, 0)
}

func newSpaciousListDelegate(selectionColor lipgloss.Color, showDesc bool, spacing int) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = showDesc
	d.SetSpacing(spacing)

	d.Styles.NormalTitle = d.Styles.
		NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#fbf1c7"})

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Border(lipgloss.Border{Left: "▎"}, false, false, false, true).
		Foreground(selectionColor).
		BorderLeftForeground(selectionColor)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	d.Styles.FilterMatch = lipgloss.NewStyle()

	return d
}
