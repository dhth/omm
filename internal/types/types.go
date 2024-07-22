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
	compactPrefixPadding  = 20
	spaciousPrefixPadding = 80
	createdAtPadding      = 40
	GOOSDarwin            = "darwin"
	taskSummaryWidth      = 120
	TaskSummaryMaxLen     = 300
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

	TaskPrefixEmptyErr      = errors.New("Task prefix cannot be empty")
	TaskSummaryBodyEmptyErr = errors.New("Task summary body is empty")
	TaskSummaryTooLongErr   = fmt.Errorf("Task summary is too long; max length allowed: %d", TaskSummaryMaxLen)
)

type Task struct {
	ItemTitle string
	ID        uint64
	Summary   string
	Context   *string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
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
		return false, TaskPrefixEmptyErr
	}

	summEls := strings.Split(summary, PrefixDelimiter)
	if len(summEls) > 1 {
		if strings.TrimSpace(summEls[0]) == "" {
			return false, TaskPrefixEmptyErr
		}

		if strings.TrimSpace(strings.Join(summEls[1:], PrefixDelimiter)) == "" {
			return false, TaskSummaryBodyEmptyErr
		}
	}

	return true, nil
}

func (t *Task) SetTitle(compact bool) {
	summEls := strings.Split(t.Summary, PrefixDelimiter)

	if compact {
		var summ string
		if len(summEls) > 1 {
			prefix := utils.RightPadTrim(summEls[0], compactPrefixPadding, true)
			summ = prefix + strings.Join(summEls[1:], PrefixDelimiter)
		} else {
			summ = t.Summary
		}

		var hasContext string
		if t.Context != nil {
			hasContext = "(c)"
		}
		t.ItemTitle = fmt.Sprintf("%s%s", utils.RightPadTrim(summ, taskSummaryWidth, true), hasContext)
		return
	}

	if len(summEls) == 1 {
		t.ItemTitle = t.Summary
		return
	}

	t.ItemTitle = utils.Trim(strings.TrimSpace(strings.Join(summEls[1:], PrefixDelimiter)), taskSummaryWidth)
}

func (t Task) Title() string {
	return t.ItemTitle
}

func (t Task) Description() string {
	var prefix string
	var createdAt string
	var hasContext string

	summEls := strings.Split(t.Summary, PrefixDelimiter)
	if len(summEls) > 1 {
		prefix = getDynamicStyle(summEls[0]).Render(utils.RightPadTrim(summEls[0], spaciousPrefixPadding, true))
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

func (p TaskPrefix) Title() string {
	return string(p)
}

func (p TaskPrefix) Description() string {
	return ""
}

func (p TaskPrefix) FilterValue() string {
	return string(p)
}
