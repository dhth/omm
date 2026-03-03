package cmd

import (
	"database/sql"
	"embed"
	"fmt"
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

	guideEntries = []entry{
		{summary: "guide: welcome to omm", filename: "01-guide-welcome-to-omm.md"},
		{summary: "domain: tasks", filename: "02-domain-tasks.md"},
		{summary: "domain: task state", filename: "03-domain-task-state.md"},
		{summary: "domain: an archived task", filename: "04-domain-an-archived-task.md", archived: true},
		{summary: "domain: task details", filename: "05-domain-task-details.md"},
		{summary: "visuals: list density", filename: "06-visuals-list-density.md"},
		{summary: "visuals: toggling context pane", filename: "07-visuals-toggling-context-pane.md"},
		{summary: "visuals: switching themes", filename: "08-visuals-switching-themes.md"},
		{summary: "actions: adding tasks", filename: "09-actions-adding-tasks.md"},
		{summary: "actions: choosing a prefix", filename: "10-actions-choosing-a-prefix.md"},
		{summary: "cli: adding a task via the CLI", filename: "11-cli-adding-a-task-via-the-cli.md"},
		{summary: "cli: importing several tasks via the CLI", filename: "12-cli-importing-several-tasks-via-the-cli.md"},
		{summary: "actions: adding context", filename: "13-actions-adding-context.md"},
		{summary: "actions: markdown in context", filename: "14-actions-markdown-in-context.md"},
		{summary: "actions: filtering tasks", filename: "15-actions-filtering-tasks.md"},
		{summary: "actions: duplicating tasks", filename: "16-actions-duplicating-tasks.md"},
		{summary: "actions: quick filtering via a list", filename: "17-actions-quick-filtering-via-a-list.md"},
		{summary: "actions: deleting a task", filename: "18-actions-deleting-a-task.md"},
		{summary: "domain: task bookmarks", filename: "19-domain-task-bookmarks.md"},
		{summary: "domain: task priorities", filename: "20-domain-task-priorities.md"},
		{summary: "actions: updating task details", filename: "21-actions-updating-task-details.md"},
		{summary: "config: changing the defaults", filename: "22-config-changing-the-defaults.md"},
		{summary: "config: flags, env vars, and config file", filename: "23-config-flags-env-vars-and-config-file.md"},
		{summary: "config: a sample TOML config", filename: "24-config-a-sample-toml-config.md"},
		{summary: "guide: and that's it!", filename: "25-guide-and-thats-it.md"},
	}
)

type entry struct {
	summary  string
	filename string
	archived bool
}

func insertGuideTasks(db *sql.DB) error {
	tasks := make([]types.Task, 0, len(guideEntries))

	now := time.Now()

	ctxs := make([]string, 0, len(guideEntries))

	for _, e := range guideEntries {
		ctx, err := getContext(e.filename)
		if err != nil {
			continue
		}
		ctxs = append(ctxs, ctx)

		tasks = append(tasks, types.Task{
			Summary:   e.summary,
			Context:   &ctxs[len(ctxs)-1],
			Active:    !e.archived,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	_, err := pers.InsertTasks(db, tasks, true)

	return err
}

func getContext(filename string) (string, error) {
	ctxBytes, err := guideFolder.ReadFile(fmt.Sprintf("%s/%s", guideAssetsPathPrefix, filename))
	if err != nil {
		return "", err
	}

	return string(ctxBytes), nil
}
