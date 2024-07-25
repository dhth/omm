package ui

import (
	"fmt"
	"strings"

	"github.com/dhth/omm/internal/types"
)

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
