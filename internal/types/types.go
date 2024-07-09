package types

import "time"

type Task struct {
	ID        uint64
	Summary   string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Task) FilterValue() string { return t.Summary }
