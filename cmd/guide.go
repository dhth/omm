package cmd

import (
	"database/sql"
	"fmt"
	"time"

	pers "github.com/dhth/omm/internal/persistence"
	"github.com/dhth/omm/internal/types"
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
			`Hi there 👋 Thanks for trying out omm.

This is a guided walkthrough to get you acquainted with omm's features.

Before we begin, let's get the basics out of the way: you exit omm by pressing
q/esc/ctrl+c. These keys also move you back menus/panes whilst using omm's TUI.

Onwards with the walkthrough then! Simply press j/↓, and follow the
instructions.
`,
			true,
		},
		{
			"guide: tasks",
			`omm is a task manager. You can also think of it as a keyboard driven to-do list.

As such, tasks are at the core of omm. A task can be anything that you want to
keep track of, ideally something that is concrete and has a clear definition of
done.

A task has a one liner summary, and optionally, some context associated with it
(like this paragraph). You can choose to add context to a task when you want to
save details that don't fit a single line.

You can view details for a task by pressing "d". Try it out (come back to this
pane by pressing q/esc/ctrl+c).

You can also update a task's summary by pressing "u".
`,
			true,
		},
		{
			"guide: task state",
			`A task can be in one of two states: active or archived.

This list shows active tasks.

To be pedantic about things, only the tasks in active list are supposed to be
"on your mind". However, there are benefits to having a list of archived lists
as well.

Press <tab> to see the archived list.
`,
			true,
		},
		{
			"guide: an archived task",
			`This is the archived list, meaning it holds tasks that are no longer being
worked on.

omm provides this list both for historical reference, and for you to be able to
move an archived task back into the active list.

You can toggle the state of a task using ctrl+d.

Press <tab/q/esc/ctrl+c> to go back to the active list.
`,
			false,
		},
		{
			"guide: list density",
			`omm's task lists can be viewed in two density modes: compact and spacious.

This is the compact mode. As opposed to this, the spacious mode shows tasks in a
more roomier list, alongside highlighting prefixes, and showing creation
timestamps. Since the list in this mode takes more space, the context pane is
shorter than the one in the compact mode. Choose whichever mode fits your
workflow better.

You can toggle between the two by pressing v. Try it now. Come back to this mode
once you're done.
`,
			true,
		},
		{
			"guide: toggling context pane",
			fmt.Sprintf(`The context pane can be toggled on/off using the backtick key("%s").

You can choose to display it or not based on your preference. For convenience,
the list will always highlight tasks that have a context associated with them by
having a "(c)" marker on them.

You can also start omm with the context pane hidden by using the flag
"--context-pane=false".
`, "`"),
			true,
		},
		{
			"guide: adding tasks",
			`Let's get to the crux of omm: adding and prioritizing tasks.

We'll begin with adding tasks. You can add a task below the cursor by pressing
"a".

Once you get acquainted with omm, you'll want to have more control on the
position of the newly added task. omm offers the following keymaps for that.

  o/a            add task below cursor
  O              add task above cursor
  I              add task at the top
  A              add task at the end

Go ahead, create a task, then move to the next guided item.
`,
			true,
		},
		{
			"guide: adding context",
			`As mentioned before, once a task is created, you might want to add context to
it.

You do that by pressing "c". This will open up the text editor you've configured
via the environment variables $EDITOR or $VISUAL (looked up in that order). You
can override this behavior by passing the "editor" flag to omm, like
"--editor='vi -u NONE'". If none of these are configured, omm falls back to "vi".

Go ahead, press "c". Try changing the text, and then save the file. This context
text should get updated accordingly.
`,
			true,
		},
		{
			"guide: task priorities",
			`At its core, omm is a list that maintains a sequence of tasks based on the
priorities you assign them.

And, as we all know, priorities often change. You're probably juggling multiple
tasks on any given day. As such, omm allows you to move tasks around in the
priority order. It has the following keymaps to achieve this:

  ⏎              move task to the top
  [2-9]          move task at index [x] to top (only in compact view)
  J              move task one position down
  K              move task one position up

It's recommended that you move the task that you're currently focussing on to
the top.
`,
			true,
		},
		{
			"guide: and that's it!",
			`That's it for the walkthrough!

I hope omm proves to be a useful tool for your task management needs. If you
find any bugs in it, or have feature requests, feel free to submit them at
https://github.com/dhth/omm/issues.

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
