package ui

import "fmt"

var (
	helpStr = fmt.Sprintf(`%s%s

%s

%s

%s

%s
%s

%s
%s

%s
%s

%s
%s

%s
%s
`, helpHeadingStyle.Render("omm"),
		helpSectionStyle.Render(` ("on-my-mind") is a keyboard-driven task manager for the command line.

Tip: Run "omm guide" for a guided walkthrough of omm's features.`),
		helpHeadingStyle.Render("omm has 5 components"),
		helpSectionStyle.Render(`1:    Active Tasks List
2:    Archived Tasks List
3:    Task creation/update Pane
4:    Task Details Pane
5:    Context Bookmarks List`),
		helpHeadingStyle.Render("Keymaps"),
		helpSubHeadingStyle.Render("General"),
		helpSectionStyle.Render(`q/esc/ctrl+c       go back
Q                  quit from anywhere`),
		helpSubHeadingStyle.Render("Active/Archived Tasks List"),
		helpSectionStyle.Render(`j/↓                move cursor down
k/↑                move cursor up
h                  go to previous page
l                  go to next page
g                  go to the top
G                  go to the end
tab                move between lists
C                  toggle showing context
d                  toggle Task Details pane
b                  open context bookmarks list
B                  open all bookmarks in the current task's context
c                  update context for a task
ctrl+d             archive/unarchive task
ctrl+x             delete task
v                  toggle between compact and spacious view`),
		helpSubHeadingStyle.Render("Active Tasks List"),
		helpSectionStyle.Render(`q/esc/ctrl+c       quit
o/a                add task below cursor
O                  add task above cursor
I                  add task at the top
A                  add task at the end
u                  update task summary
⏎                  move task to the top
[2-9]              move task at index [x] to top (only in compact view)
J                  move task one position down
K                  move task one position up`),
		helpSubHeadingStyle.Render("Task Details Pane"),
		helpSectionStyle.Render(`h/l                move backwards/forwards when in the task details view
B                  open all bookmarks in the current task's context`),
		helpSubHeadingStyle.Render("Context Bookmarks List"),
		helpSectionStyle.Render(`⏎                  open URL in browser`),
	)
)
