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
	guideAssetsPathPrefix = "assets/guide"
)

var (
	//go:embed assets/guide/*.md
	guideFolder embed.FS

	guideSummaryRegex = regexp.MustCompile(`[^a-z-]`)

	guideEntries = []entry{
		{summary: "guide: welcome to omm"},
		{summary: "domain: tasks"},
		{summary: "domain: task state"},
		{summary: "domain: an archived task", archived: true},
		{summary: "domain: task details"},
		{summary: "visuals: list density"},
		{summary: "visuals: toggling context pane"},
		{summary: "actions: adding tasks"},
		{summary: "cli: adding a task via the CLI"},
		{summary: "cli: importing several tasks via the CLI"},
		{summary: "actions: adding context"},
		{summary: "actions: filtering tasks"},
		{summary: "actions: quick filtering via a list"},
		{summary: "domain: task bookmarks"},
		{summary: "domain: task priorities"},
		{summary: "actions: updating task details"},
		{summary: "config: changing the defaults"},
		{summary: "config: flags, env vars, and config file"},
		{summary: "config: a sample TOML config"},
		{summary: "guide: and that's it!"},
	}
)

type entry struct {
	summary  string
	archived bool
}

func getContext(summary string) (string, error) {
	summary = strings.ToLower(summary)
	summary = strings.ReplaceAll(summary, " ", "-")
	fPath := guideSummaryRegex.ReplaceAllString(summary, "")

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

		if err != nil {
			continue
		}

		tasks[i] = types.Task{
			Summary:   e.summary,
			Context:   &ctxs[i],
			Active:    !e.archived,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	err = pers.InsertTasksIntoDB(db, tasks)

	return err
}
