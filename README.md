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

**go**:

```sh
go install github.com/dhth/omm@latest
```

⚡️ Usage
---

### TUI

`omm`'s TUI allows you to do the following (all with one or two key presses):

- create/update tasks at a specific position in the priority list
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

### Importing tasks

Multiple tasks can be imported from `stdin` using the `import` subcommand.

```bash
cat << 'EOF' | omm import
Order new ACME rocket skates
Draw fake tunnel on the canyon wall
Assemble ACME jet-propelled pogo stick
EOF
```

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

Screenshots
---

![Screen 1](https://tools.dhruvs.space/images/omm/omm-1.png)

![Screen 2](https://tools.dhruvs.space/images/omm/omm-2.png)

⌨️ Keymaps
---

```text
j/↓: move cursor down                           k/↑: move cursor up
o: add task above cursor                        a: add task below cursor
O: add task at the beginning of list            A: add task at the end of list
J: move task one position down                  K: move task one position up
[2-9]: move task at index "x" to top of list    ⏎: move task on cursor to top of list
ctrl+d: archive/unarchive task                  ctrl+x: delete task
g: go to top of list                            G: go to end of list
tab: move between views                         q/esc: go back/quit
```

Acknowledgements
---

`omm` is built using [bubbletea][1], and is released using [goreleaser][2], both
of which are amazing tools.

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/goreleaser/goreleaser
