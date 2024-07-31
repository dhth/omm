package ui

import (
	"time"

	"github.com/dhth/omm/internal/types"
)

type HideHelpMsg struct{}

type taskSequenceUpdatedMsg struct {
	err error
}

type taskCreatedMsg struct {
	id          uint64
	taskSummary string
	createdAt   time.Time
	updatedAt   time.Time
	err         error
}

type taskDeletedMsg struct {
	id        uint64
	listIndex int
	active    bool
	err       error
}

type taskSummaryUpdatedMsg struct {
	listIndex   int
	id          uint64
	taskSummary string
	updatedAt   time.Time
	err         error
}

type taskContextUpdatedMsg struct {
	listIndex int
	list      taskListType
	id        uint64
	context   string
	updatedAt time.Time
	err       error
}

type taskStatusChangedMsg struct {
	listIndex int
	id        uint64
	active    bool
	updatedAt time.Time
	err       error
}

type tasksFetched struct {
	tasks  []types.Task
	active bool
	err    error
}

type textEditorClosed struct {
	fPath      string
	taskIndex  int
	taskId     uint64
	oldContext *string
	err        error
}

type uriOpenedMsg struct {
	url string
	err error
}

type urisOpenedDarwinMsg struct {
	urls []string
	err  error
}

type contextWrittenToCBMsg struct {
	err error
}
