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

omm offers guided walkthrough of its features, intended for new users of it. Run
it as follows.

```bash
omm guide
```

![Guide](https://tools.dhruvs.space/images/omm/omm-guide-1.png)

⚡️ Usage
---

### TUI

`omm`'s TUI is comprised of several panes: 2 lists (for active and archived
tasks), a context pane, and a task entry/update pane.

You can do the following via the TUI:

- create/update tasks at a specific position in the priority list
- add/update context for a task
- add a task at the start/end of the list
- move a task to the top of the list (indicating that it takes the highest
    priority at the moment)
- move task up/down based on changing priorities
- archive/unarchive a task
- view archived tasks
- delete a task

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

Compact mode:

![compact mode](https://tools.dhruvs.space/images/omm/omm-1.png)

Spacious mode:

![spacious mode](https://tools.dhruvs.space/images/omm/omm-2.png)

### Importing tasks

Multiple tasks can be imported from `stdin` using the `import` subcommand.

```bash
cat << 'EOF' | omm import
Order new ACME rocket skates
Draw fake tunnel on the canyon wall
Assemble ACME jet-propelled pogo stick
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

Outputting tasks
---

Tasks can be outputted to `stdout` using the `tasks` subcommand.

```bash
omm tasks
```

⌨️ Keymaps
---

```text
j/↓            move cursor down
k/↑            move cursor up
o/a            add task below cursor
O              add task above cursor
I              add task at the top
A              add task at the end
u              update task summary
⏎              move task to the top
[2-9]          move task at index [x] to top (only in compact view)
J              move task one position down
K              move task one position up
ctrl+d         archive/unarchive task
ctrl+x         delete task
g              go to the top
G              go to the end
tab            move between views
c              update context for a task
d              show task details in a full screen pane
v              toggle between compact and spacious view
`              toggle showing context
q/esc/ctrl+c   go back/quit
```

Acknowledgements
---

`omm` is built using [bubbletea][1], and is released using [goreleaser][2], both
of which are amazing tools.

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/goreleaser/goreleaser
