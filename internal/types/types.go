package types

import (
	"errors"
	"strings"
	"time"
)

const (
	PrefixDelimiter   = ":"
	GOOSDarwin        = "darwin"
	TaskSummaryMaxLen = 300
)

var (
	ErrTaskSummaryEmpty     = errors.New("task summary is empty")
	ErrTaskPrefixEmpty      = errors.New("task prefix is empty")
	ErrTaskSummaryBodyEmpty = errors.New("task summary body is empty")
	ErrTaskSummaryTooLong   = errors.New("task summary is too long")
)

type TaskDetails struct {
	Summary string
	Context *string
}

type Task struct {
	ID        uint64
	Summary   string
	Context   *string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
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
	// description is handled by custom delegates (compactItemDelegate/spaciousTaskItemDelegate)
	return ""
}

func (t Task) FilterValue() string {
	p, ok := t.Prefix()
	if ok {
		return string(p)
	}
	return ""
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
