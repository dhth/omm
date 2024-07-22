package cmd

import (
	"database/sql"
	"fmt"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
)

const (
	codeBlockMarker = "```"
)

type entry struct {
	summary string
	ctx     string
	active  bool
}

func insertGuideTasks(db *sql.DB) error {
	var entries = []entry{
		{
			"guide: welcome to omm",
			fmt.Sprintf(`Hi there 👋 Thanks for trying out **omm**.

This is a guided walkthrough to get you acquainted with omm's features.

Before we begin, let's get the basics out of the way: you exit omm by pressing %s. These keys also move you back menus/panes whilst using omm's TUI.

Onwards with the walkthrough then! Simply press %s, and follow the instructions.
`, "`q/esc/<ctrl+c>`", "`j/↓`"),
			true,
		},
		{
			"domain: tasks",
			`omm (**on-my-mind**) is a task manager. You can also think of it as a keyboard driven to-do list.

As such, tasks are at the core of omm. A task can be anything that you want to keep track of, ideally something that is concrete and has a clear definition of done.

Tasks in omm have a one liner summary, and optionally, some context associated with them (like this paragraph). You can choose to add context to a task when you want to save details that don't fit in a single line.
`,
			true,
		},
		{
			"domain: task state",
			fmt.Sprintf(`A task can be in one of two states: **active** or **archived**.

This list shows active tasks.

To be pedantic about things, only the tasks in the active list are supposed to be "on your mind". However, there are benefits to having a list of archived tasks as well.

Press %s to see the archived list.
`, "`<tab>`"),
			true,
		},
		{
			"domain: task details",
			fmt.Sprintf(`The **Task Details** pane is intended for when you simply want to read all the details associated with a task in a full screen view.

You can view this pane by pressing %s. This pane is useful when a task's context is too long to fit in the context pane.

Whilst in this pane, you can move backwards and forwards in the task list by pressing %s. You quit out of this pane by either pressing %s again, or %s.

Try it out. Come back to this entry when you're done.
`, "`d`", "`h/←/→/l`", "`d`", "`q/esc/<ctrl+c>`"),
			true,
		},
		{
			"domain: an archived task",
			fmt.Sprintf(`This is the **Archived Tasks** list, meaning it holds tasks that are no longer being worked on.

omm provides this list both for historical reference, as well as for you to be able to move an archived task back into the active list.

You can toggle the state of a task using %s.

Press %s to go back to the active list.
`, "`ctrl+d`", "`tab/q/esc/<ctrl+c>`"),
			false,
		},
		{
			"visuals: list density",
			fmt.Sprintf(`omm's task lists can be viewed in two density modes: **compact** and **spacious**.

This is the compact mode. As opposed to this, the spacious mode shows tasks in a more roomier list, alongside highlighting prefixes (we'll see what that means), and showing creation timestamps.

omm starts up with spacious mode by default (you can change this, as we'll see soon). You can toggle between the two modes by pressing %s. Choose whichever mode fits your workflow better.

Try it out. Come back to this mode once you're done.
`, "`v`"),
			true,
		},
		{
			"visuals: toggling context pane",
			fmt.Sprintf(`The context pane can be toggled on/off by pressing %s.

You can choose to display it or not based on your preference. For convenience, the lists will always highlight tasks that have a context associated with them by having a **(c)** marker on them.

omm starts up with the context pane hidden by default (you can change this, as we'll see soon).
`, "`C`"),
			true,
		},
		{
			"actions: adding tasks",
			fmt.Sprintf(`Let's get to the crux of omm: **adding** and **prioritizing** tasks. We'll begin with adding tasks.

You can add a task below the cursor by pressing %s. Once you get acquainted with omm, you'll want to have more control on the position of the newly added task. omm offers the following keymaps for that.

    o/a            add task below cursor
    O              add task above cursor
    I              add task at the top
    A              add task at the end

Go ahead, create a task, then move to the next guided item.
`, "`a`"),
			true,
		},
		{
			"cli: adding a task via the CLI",
			fmt.Sprintf(`You can also add a task to omm via its command line interface. For example:

%sbash
omm 'prefix: a task summary'
%s

This will add an entry to the top of the active tasks list.
`, codeBlockMarker, codeBlockMarker),
			true,
		},
		{
			"cli: importing several tasks via the CLI",
			fmt.Sprintf(`You can also import more than one task at a time by using the **import** subcommand. For example:

%sbash
cat << 'EOF' | omm import
orders: order new ACME rocket skates
traps: draw fake tunnel on the canyon wall
tech: assemble ACME jet-propelled pogo stick
EOF
%s

omm will expect each line in stdin to hold one task's summary.
`, codeBlockMarker, codeBlockMarker),
			true,
		},
		{
			"actions: adding context",
			fmt.Sprintf(`As mentioned before, once a task is created, you might want to add context to it.

You do that by pressing %s. Go ahead, try it out. Try changing the text, and then save the file. This context text should get updated accordingly.

Once saved, you can also copy a tasks's context to your system clipboard by pressing %s.
`, "`c`", "`y`"),
			true,
		},
		{
			"actions: filtering tasks",
			fmt.Sprintf(`You can filter tasks in a list by pressing %s. Doing this will open up a search prompt, which will match your query with task prefixes.

Try it out now. You get out of the filtered state by pressing %s.

Note: You cannot add tasks or move them around in a filtered state. But, you can move a task to the top of the list (by pressing %s). Doing this will also get you out of the filtered state.
`, "`/`", "`q/esc/<ctrl+c>`", "`⏎`"),
			true,
		},
		{
			"actions: quick filtering via a list",
			fmt.Sprintf(`You can also choose the prefix you want to filter by with the means of a list, hereby called as the **Quick Filter List**. Press %s to open up a set of task prefixes contained in the currently active task list. Press %s to pre-populate the task list's search prompt with your selection.

Try it out now.

Note: Both the **Active Tasks List** and **Archived Tasks List** can be filtered separately, using either the manual filtering approach or via the **Quick Filter List**.
`, "`ctrl+p`", "`⏎`"),
			true,
		},
		{
			"domain: task bookmarks",
			fmt.Sprintf(`Sometimes you'll add URLs to a task's summary or its context.

Such URLs (eg. https://github.com/dhth/omm, https://tools.dhruvs.space, https://c.xkcd.com/random/comic) could be placed anywhere in the summary/context.

omm lets you open these URLs via a single keypress. You can either press %s to open up a list of all URLs, and then open one of them by pressing %s, or open all of them by pressing %s.

Note: If a task has a single URL added to it, pressing %s will skip showing the list, and open the URL directly.

Try both approaches now.
`, "`b`", "`⏎`", "`B`", "`b`"),
			true,
		},
		{
			"domain: task priorities",
			`At its core, omm is a dynamic list that maintains a sequence of tasks based on the priorities you assign them.

And, as we all know, priorities often change. You're probably juggling multiple tasks on any given day. As such, omm allows you to move tasks around in the priority order. It has the following keymaps to achieve this:

    ⏎              move task to the top
    J              move task one position down
    K              move task one position up

It's recommended that you move the task that you're currently focussing on to the top.
`,
			true,
		},
		{
			"actions: updating task details",
			fmt.Sprintf(`Once a task is created, its summary and context can be changed at any point.

You can update a task's summary by pressing %s.

This will open up the the same prompt you saw when creating a new task, with the only difference that the task's summary will be pre-filled for you. This can come in handy when you want to quickly jot down a task for yourself (either by using the TUI, or by using the CLI (eg. %s)), and then come back to it later to refine it more.

Similarly, you can also update a task's context any time (by pressing %s).
`, "`u`", "`omm 'a hastily written task summary'`", "`c`"),
			true,
		},
		{
			"config: changing the defaults",
			fmt.Sprintf(`omm allows you to change the some of its behavior via configuration, which it will consider in the order listed below:

- CLI flags (run %s to see details)
- Environment variables (eg. %s)
- A TOML configuration file (run %s to see where this lives; you can change this via the flag %s)

omm will consider configuration in the order laid out above, ie, CLI flags will take the highest priority.
`, "`omm -h`", "`$OMM_EDITOR`", "`omm -h`", "`--config-path`"),
			true,
		},
		{
			"config: flags, env vars, and config file",
			fmt.Sprintf(`Every flag listed by %s (except %s) has an environment variable counterpart, as well as a TOML config counterpart. 

For example:

%stext
--show-context  ->  OMM_SHOW_CONTEXT  ->  show_context
--editor        ->  OMM_EDITOR        ->  editor
%s
`, "`omm -h`", "`--config-path`", codeBlockMarker, codeBlockMarker),
			true,
		},
		{
			"config: a sample TOML config",
			fmt.Sprintf(`Here's a sample TOML configuration file:

%stoml
%s
%s
`, codeBlockMarker, sampleCfg, codeBlockMarker),
			true,
		},
		{
			"guide: and that's it!",
			`That's it for the walkthrough!

I hope omm proves to be a useful tool for your task management needs. If you find any bugs in it, or have feature requests, feel free to submit them at https://github.com/dhth/omm/issues.

Happy task managing! 👋
`,
			true,
		},
	}

	tasks := make([]types.Task, len(entries))

	now := time.Now()
	for i, e := range entries {
		tasks[i] = types.Task{
			Summary:   e.summary,
			Context:   &e.ctx,
			Active:    e.active,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	err := pers.InsertTasksIntoDB(db, tasks)

	return err
}
