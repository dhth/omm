package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/types"
	"github.com/dhth/omm/internal/utils"
)

type compactItemDelegate struct {
	selStyle lipgloss.Style
}

func (d compactItemDelegate) Height() int { return 1 }

func (d compactItemDelegate) Spacing() int { return 1 }

func (d compactItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d compactItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(types.Task)
	if !ok {
		return
	}

	prefix, sc, hp := t.GetPrefixAndSummaryContent()

	if hp {
		prefix = types.GetDynamicStyle(prefix).Render(utils.RightPadTrim(prefix, prefixPadding, true))
	}
	var hasContext string
	if t.Context != nil {
		hasContext = "(c)"
	}

	sr := d.selStyle.Render
	var str string
	if index == m.Index() {
		str = fmt.Sprintf("%s%s%s%s", sr("│ "), prefix, sr(utils.RightPadTrim(sc, taskSummaryWidth-prefixPadding, true)), sr(hasContext))
	} else {
		str = fmt.Sprintf("%s%s%s%s", "  ", prefix, utils.RightPadTrim(sc, taskSummaryWidth-prefixPadding, true), hasContext)
	}

	fmt.Fprint(w, str)
}

func newSpaciousListDelegate(color lipgloss.Color, showDesc bool, spacing int) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = showDesc
	d.SetSpacing(spacing)

	d.Styles.NormalTitle = d.Styles.
		NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#fbf1c7"})

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	d.Styles.FilterMatch = lipgloss.NewStyle()

	return d
}
