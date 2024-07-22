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

%s
%s
`, helpHeadingStyle.Render("omm"),
		helpSectionStyle.Render(` ("on-my-mind") is a keyboard-driven task manager for the command line.

Tip: Run "omm guide" for a guided walkthrough of omm's features.`),
		helpHeadingStyle.Render("omm has 6 components"),
		helpSectionStyle.Render(`1:    Active Tasks List
2:    Archived Tasks List
3:    Task creation/update Pane
4:    Task Details Pane
5:    Task Bookmarks List
6:    Quick Filter List`),
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
b                  open Task Bookmarks list
B                  open all bookmarks added to current task
c                  update context for a task
ctrl+d             archive/unarchive task
ctrl+x             delete task
ctrl+r             reload task lists
/                  filter list by task prefix
ctrl+p             open the "Quick Filter List"
y                  copy selected task's context to system clipboard
v                  toggle between compact and spacious view`),
		helpSubHeadingStyle.Render("Active Tasks List"),
		helpSectionStyle.Render(`q/esc/ctrl+c       quit
o/a                add task below cursor
O                  add task above cursor
I                  add task at the top
A                  add task at the end
u                  update task summary
⏎                  move task to the top
J                  move task one position down
K                  move task one position up`),
		helpSubHeadingStyle.Render("Task Details Pane"),
		helpSectionStyle.Render(`h/l                move backwards/forwards when in the task details view
y                  copy current task's context to system clipboard
B                  open all bookmarks added to current task`),
		helpSubHeadingStyle.Render("Task Bookmarks List"),
		helpSectionStyle.Render(`⏎                  open URL in browser`),
		helpSubHeadingStyle.Render("Quick Filter List"),
		helpSectionStyle.Render(`⏎                  pre-populate task list's search prompt with chosen prefix`),
	)
)
