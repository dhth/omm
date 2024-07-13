package ui

import "fmt"

var helpStr = fmt.Sprintf(`omm ("on-my-mind") is a keyboard-driven task manager for the command line.

[Keymaps]
j/↓                move cursor down
k/↑                move cursor up
o/a                add task below cursor
O                  add task above cursor
I                  add task at the top
A                  add task at the end
u                  update task summary
⏎                  move task to the top
[2-9]              move task at index [x] to top (only in compact view)
J                  move task one position down
K                  move task one position up
ctrl+d             archive/unarchive task
ctrl+x             delete task
g                  go to the top
G                  go to the end
tab                move between views
c                  update context for a task
d                  show task details in a full screen pane
v                  toggle between compact and spacious view
%s                  toggle showing context
q/esc/ctrl+c       go back/quit

Run "omm guide" for a guided walkthrough of omm's features.
`, "`")
