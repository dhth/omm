# omm

✨ Overview
---

`omm` (stands for "on-my-mind") is a keyboard-driven task manager for the
command line.

`omm` is intended for those who need to frequently rearrange tasks in their
to-do list. It lets you move any item to the top of the list, add tasks at
specific positions, and adjust task priorities up or down, all via one or two
keypresses.

![Usage](https://tools.dhruvs.space/images/omm/omm.gif)

[source video](https://www.youtube.com/watch?v=_VnvgqVdU20)

🤔 Motivation
---

The fundamental idea behind `omm` is that while we might have several tasks on
our to-do list — each with its own priority — we typically focus on one task at
a time. Priorities frequently change, requiring us to switch between tasks.
`omm` lets you visualize this shifting priority order with a very simple list
interface that can be managed entirely via the keyboard.

💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/omm
```

**go**:

```sh
go install github.com/dhth/omm@latest
```

Or get the binaries directly from a
[release](https://github.com/dhth/omm/releases).

💡 Guide
---

omm offers a guided walkthrough of its features, intended for new users of it.
Run it as follows.

```bash
omm guide
```

![Guide](https://tools.dhruvs.space/images/omm/omm-guide-1.png)

⚡️ Usage
---

### TUI

`omm`'s TUI is comprised of several panes: 3 lists (for active and archived
tasks, and one for task bookmarks), a context pane, and a task entry/update
pane.

#### Active Tasks List

As the name suggests, the active tasks list is for the tasks you're actively
working on right now. It allows you to do the following:

- Create/update tasks at a specific position in the priority list
- Add a task at the start/end of the list
- Move a task to the top of the list (indicating that it takes the highest
    priority at the moment)
- Move task up/down based on changing priorities
- Archive a task
- Permanently delete a task

![active-tasks](https://tools.dhruvs.space/images/omm/omm-active-tasks-1.png)

#### Archived Tasks List

Once you're done with a task, you can archive it, which puts it in the archived
tasks list. It's more for historical reference, but you can also unarchive a
task and put it back in the active list, if you need to. You can also
permanently delete tasks from here.

![active-tasks](https://tools.dhruvs.space/images/omm/omm-archived-tasks-1.png)

#### Context Pane

For tasks that need more details that you can fit in a one line summary, there
is the context pane. You add/update context for a task via a text editor which
is chosen based on the following look ups:

- the "--editor" flag
- $OMM_EDITOR
- "editor" property in omm's toml config
- $EDITOR/$VISUAL
- `vi` (fallback)

![active-tasks](https://tools.dhruvs.space/images/omm/omm-context-1.png)

**[`^ back to top ^`](#omm)**

#### Task Entry Pane

This is where you enter/update a task summary. If you enter a summary in the
format `prefix: task summary goes here`, `omm` will highlight the prefix for you
in the task lists.

![active-tasks](https://tools.dhruvs.space/images/omm/omm-task-entry-1.png)

#### Tweaking the TUI

The list colors and the task list title can be changed via CLI flags.

```bash
omm \
    --tl-color="#b8bb26" \
    --atl-color="#fb4934" \
    --title="work"
```

omm offers two modes for the visual density of its lists: "compact" and
"spacious", the former being the default. omm can be started with one of
the two modes, which can later be switched by pressing "v".

```bash
omm --list-density=spacious
```

This configuration property can also be provided via the environment variable
`OMM_LIST_DENSITY`.

Compact mode:

![compact](https://tools.dhruvs.space/images/omm/omm-compact-1.png)

Spacious mode:

![spacious](https://tools.dhruvs.space/images/omm/omm-spacious-1.png)

### Importing tasks

Multiple tasks can be imported from `stdin` using the `import` subcommand.

```bash
cat << 'EOF' | omm import
orders: order new ACME rocket skates
traps: draw fake tunnel on the canyon wall
tech: assemble ACME jet-propelled pogo stick
EOF
```

Tip: Vim users can import tasks into omm by making a visual selection and
running `:'<,'>!omm import<CR>`.

### Adding a single task

When an argument is passed to `omm`, it saves it as a task, instead of opening
up the TUI.

```bash
omm "Install spring-loaded boxing glove"
```

### Configuration

`omm` allows you to change the some of its behavior via configuration, which it
will consider in the order listed below:

- CLI flags (run `omm -h` to see details)
- Environment variables (eg. `OMM_EDITOR`)
- A TOML configuration file (run `omm -h` to see where this lives; you can
    change this via the flag `--config-path`)

    Here's a sample config file:

    ```toml
    db_path      = "~/.local/share/omm/omm-w.db"
    tl_color     = "#b8bb26"
    atl_color    = "#fabd2f"
    title        = "work"
    list_density = "spacious"
    show_context = false
    editor       = "vi -u NONE"
    ```

**[`^ back to top ^`](#omm)**

Outputting tasks
---

Tasks can be outputted to `stdout` using the `tasks` subcommand.

```bash
omm tasks
```

⌨️ Keymaps
---

```text
General
q/esc/ctrl+c       go back
Q                  quit from anywhere

Active/Archived Tasks List
j/↓                move cursor down
k/↑                move cursor up
h                  go to previous page
l                  go to next page
g                  go to the top
G                  go to the end
tab                move between lists
C                  toggle showing context
d                  toggle Task Details pane
b                  open task bookmarks list
B                  open all bookmarks added to current task
c                  update context for a task
ctrl+d             archive/unarchive task
ctrl+x             delete task
ctrl+r             reload task lists
y                  copy selected task's context to system clipboard
v                  toggle between compact and spacious view

Active Tasks List
q/esc/ctrl+c       quit
o/a                add task below cursor
O                  add task above cursor
I                  add task at the top
A                  add task at the end
u                  update task summary
⏎                  move task to the top
[2-9]              move task at index [x] to top (only in compact view)
J                  move task one position down
K                  move task one position up

Task Details Pane
h/l                move backwards/forwards when in the task details view
y                  copy selected task's context to system clipboard
B                  open all bookmarks added to current task

Task Bookmarks List
⏎                  open URL in browser
```

Acknowledgements
---

`omm` stands on the shoulders of giants. 

- [bubbletea](https://github.com/charmbracelet/bubbletea) as the TUI framework
- [sqlite](https://www.sqlite.org) as the local database
- [goreleaser](https://github.com/goreleaser/goreleaser) for releasing binaries

**[`^ back to top ^`](#omm)**
