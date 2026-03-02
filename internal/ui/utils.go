package ui

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/types"
)

func getColorForString(str string, colors []string) lipgloss.TerminalColor {
	if len(colors) == 0 {
		return lipgloss.NoColor{}
	}

	h := fnv.New32()
	h.Write([]byte(str))
	hash := h.Sum32()

	color := colors[hash%uint32(len(colors))]

	return lipgloss.Color(color)
}

func getSummaryWithNewPrefix(summary, newPrefix string) string {
	if strings.TrimSpace(summary) == "" {
		return fmt.Sprintf("%s: ", newPrefix)
	}

	summEls := strings.Split(summary, types.PrefixDelimiter)

	if len(summEls) == 1 {
		return fmt.Sprintf("%s: %s", newPrefix, summary)
	}

	if summEls[1] == "" {
		return fmt.Sprintf("%s: ", newPrefix)
	}

	return fmt.Sprintf("%s:%s", newPrefix, strings.Join(summEls[1:], types.PrefixDelimiter))
}

func getPrefix(summary string) (string, bool) {
	if strings.TrimSpace(summary) == "" {
		return "", false
	}

	summEls := strings.Split(summary, types.PrefixDelimiter)

	if len(summEls) == 1 {
		return "", false
	}

	return strings.TrimSpace(summEls[0]), true
}
