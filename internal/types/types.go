package types

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/omm/internal/utils"
	"github.com/dustin/go-humanize"
)

const (
	timeFormat            = "2006/01/02 15:04"
	PrefixDelimiter       = ":"
	compactPrefixPadding  = 24
	spaciousPrefixPadding = 80
	createdAtPadding      = 40
	GOOSDarwin            = "darwin"
	taskSummaryWidth      = 120
	TaskSummaryMaxLen     = 300
)

var (
	createdAtColor  = "#928374"
	hasContextColor = "#928374"
	createdAtStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color(createdAtColor))

	hasContextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(hasContextColor))

	ErrTaskSummaryEmpty     = errors.New("task summary is empty")
	ErrTaskPrefixEmpty      = errors.New("task prefix is empty")
	ErrTaskSummaryBodyEmpty = errors.New("task summary body is empty")
	ErrTaskSummaryTooLong   = errors.New("task summary is too long")
)

type TaskDetails struct {
	Summary string  `json:"summary"`
	Context *string `json:"context"`
}

type Task struct {
	ID        uint64    `json:"-"`
	Summary   string    `json:"summary"`
	Context   *string   `json:"context"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (t Task) GetDetails() TaskDetails {
	var context *string
	if t.Context != nil {
		c := *t.Context
		context = &c
	}

	return TaskDetails{
		Summary: t.Summary,
		Context: context,
	}
}

type ContextBookmark string

type TaskPrefix string

func (t Task) Prefix() (TaskPrefix, bool) {
	summEls := strings.Split(t.Summary, PrefixDelimiter)
	if len(summEls) > 1 {
		// This shouldn't happen, but it's still good to check this to ensure
		// the quick filter list doesn't misbehave
		if strings.TrimSpace(summEls[0]) == "" {
			return "", false
		}
		return TaskPrefix(strings.TrimSpace(summEls[0])), true
	}
	return "", false
}

func CheckIfTaskSummaryValid(summary string) (bool, error) {
	if strings.TrimSpace(summary) == "" {
		return false, ErrTaskSummaryEmpty
	}

	if len(summary) > TaskSummaryMaxLen {
		return false, ErrTaskSummaryTooLong
	}

	summEls := strings.Split(summary, PrefixDelimiter)
	if len(summEls) > 1 {
		if strings.TrimSpace(summEls[0]) == "" {
			return false, ErrTaskPrefixEmpty
		}

		if strings.TrimSpace(strings.Join(summEls[1:], PrefixDelimiter)) == "" {
			return false, ErrTaskSummaryBodyEmpty
		}
	}

	return true, nil
}

func (t Task) GetPrefixAndSummaryContent() (string, string, bool) {
	summEls := strings.Split(t.Summary, PrefixDelimiter)

	if len(summEls) == 1 {
		return "", t.Summary, false
	}

	return strings.TrimSpace(summEls[0]), strings.TrimSpace(strings.Join(summEls[1:], PrefixDelimiter)), true
}

func (t Task) Title() string {
	_, sc, _ := t.GetPrefixAndSummaryContent()
	return sc
}

func (t Task) Description() string {
	var prefix string
	var createdAt string
	var hasContext string

	summEls := strings.Split(t.Summary, PrefixDelimiter)
	if len(summEls) > 1 {
		prefix = GetDynamicStyle(summEls[0]).Render(utils.RightPadTrim(summEls[0], spaciousPrefixPadding, true))
	} else {
		prefix = strings.Repeat(" ", spaciousPrefixPadding)
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

func (t Task) FilterValue() string {
	p, ok := t.Prefix()
	if ok {
		return string(p)
	}
	return ""
}

func GetDynamicStyle(str string) lipgloss.Style {
	h := fnv.New32()
	h.Write([]byte(str))
	hash := h.Sum32()

	color := colors[hash%uint32(len(colors))]
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

func (p TaskPrefix) Title() string {
	return string(p)
}

func (p TaskPrefix) Description() string {
	return ""
}

func (p TaskPrefix) FilterValue() string {
	return string(p)
}

type TaskStatusFilter uint8

const (
	TaskStatusActive TaskStatusFilter = iota
	TaskStatusInactive
	TaskStatusAny
)

func TaskStatusFilterValues() []string {
	return []string{"active", "inactive", "any"}
}

func ParseTaskStatusFilter(value string) (TaskStatusFilter, bool) {
	switch value {
	case "active":
		return TaskStatusActive, true
	case "inactive":
		return TaskStatusInactive, true
	case "any":
		return TaskStatusAny, true
	default:
		return TaskStatusAny, false
	}
}
