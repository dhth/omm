package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
)

const (
	codeBlockMarker       = "```"
	guideAssetsPathPrefix = "assets/guide"
)

var (
	//go:embed assets/guide/*.md
	guideFolder embed.FS

	guideEntries = []entry{
		{
			"guide: welcome to omm",
			true,
		},
		{
			"domain: tasks",
			true,
		},
		{
			"domain: task state",
			true,
		},
		{
			"domain: an archived task",
			false,
		},
		{
			"domain: task details",
			true,
		},
		{
			"visuals: list density",
			true,
		},
		{
			"visuals: toggling context pane",
			true,
		},
		{
			"actions: adding tasks",
			true,
		},
		{
			"cli: adding a task via the CLI",
			true,
		},
		{
			"cli: importing several tasks via the CLI",
			true,
		},
		{
			"actions: adding context",
			true,
		},
		{
			"actions: filtering tasks",
			true,
		},
		{
			"actions: quick filtering via a list",
			true,
		},
		{
			"domain: task bookmarks",
			true,
		},
		{
			"domain: task priorities",
			true,
		},
		{
			"actions: updating task details",
			true,
		},
		{
			"config: changing the defaults",
			true,
		},
		{
			"config: flags, env vars, and config file",
			true,
		},
		{
			"config: a sample TOML config",
			true,
		},
		{
			"guide: and that's it!",
			true,
		},
	}
)

type entry struct {
	summary string
	active  bool
}

func getContext(summary string) (string, error) {
	summary = strings.ToLower(summary)
	summary = strings.ReplaceAll(summary, " ", "-")
	re := regexp.MustCompile(`[^a-z-]`)
	fPath := re.ReplaceAllString(summary, "")

	ctxBytes, err := guideFolder.ReadFile(fmt.Sprintf("%s/%s.md", guideAssetsPathPrefix, fPath))
	if err != nil {
		return "", err
	}

	return string(ctxBytes), nil
}

func insertGuideTasks(db *sql.DB) error {

	tasks := make([]types.Task, len(guideEntries))

	now := time.Now()

	ctxs := make([]string, len(guideEntries))

	var err error
	for i, e := range guideEntries {
		ctxs[i], err = getContext(guideEntries[i].summary)

		if err == nil {
			tasks[i] = types.Task{
				Summary:   e.summary,
				Context:   &ctxs[i],
				Active:    e.active,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
	}

	err = pers.InsertTasksIntoDB(db, tasks)

	return err
}
