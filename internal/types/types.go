package types

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/utils"
	"github.com/dustin/go-humanize"
)

const (
	timeFormat       = "2006/01/02 15:04"
	prefixDelimiter  = ":"
	prefixPadding    = 80
	createdAtPadding = 40
	GOOSDarwin       = "darwin"
)

var (
	createdAtColor  = "#928374"
	hasContextColor = "#928374"
	taskColors      = []string{
		"#d3869b",
		"#b5e48c",
		"#90e0ef",
		"#ca7df9",
		"#ada7ff",
		"#bbd0ff",
		"#48cae4",
		"#8187dc",
		"#ffb4a2",
		"#b8bb26",
		"#ffc6ff",
		"#4895ef",
		"#83a598",
		"#fabd2f",
	}
	createdAtStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(createdAtColor))

	hasContextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(hasContextColor))
)

type Task struct {
	ID        uint64
	Summary   string
	Context   *string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContextBookmark string

func (t Task) Title() string {
	summEls := strings.Split(t.Summary, prefixDelimiter)
	if len(summEls) == 1 {
		return t.Summary
	}
	return strings.TrimSpace(strings.Join(summEls[1:], prefixDelimiter))
}

func (t Task) Description() string {
	var prefix string
	var createdAt string
	var hasContext string

	summEls := strings.Split(t.Summary, prefixDelimiter)
	if len(summEls) > 1 {
		prefix = getDynamicStyle(summEls[0]).Render(utils.RightPadTrim(summEls[0], prefixPadding, true))
	} else {
		prefix = strings.Repeat(" ", prefixPadding)
	}
	now := time.Now()

	var createdAtTs string
	if now.Sub(t.CreatedAt).Seconds() < 60 {
		createdAtTs = "just now"
	} else {
		createdAtTs = humanize.Time(t.CreatedAt)
	}
	createdAt = createdAtStyle.Render(utils.RightPadTrim(fmt.Sprintf("created %s", createdAtTs), createdAtPadding, true))

	if t.Context != nil {
		hasContext = hasContextStyle.Render("(c)")
	}

	return fmt.Sprintf("%s%s%s", prefix, createdAt, hasContext)
}

func (t Task) FilterValue() string { return t.Summary }

func getDynamicStyle(str string) lipgloss.Style {
	h := fnv.New32()
	h.Write([]byte(str))
	hash := h.Sum32()

	color := taskColors[int(hash)%len(taskColors)]
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color))
}

func (c ContextBookmark) Title() string {
	return string(c)
}

func (c ContextBookmark) Description() string {
	return ""
}

func (c ContextBookmark) FilterValue() string {
	return string(c)
}
