## Overview

**omm** ("on-my-mind") is a keyboard-driven task manager for the command line.

Tip: Run `omm guide` for a guided walkthrough of omm's features.

omm has 6 components:

- Active Tasks List
- Archived Tasks List
- Task creation/update Pane
- Task Details Pane
- Task Bookmarks List
- Quick Filter List

## Keymaps

### General

    q/esc/ctrl+c       go back
    Q                  quit from anywhere

### Active/Archived Tasks List

    j/↓                move cursor down
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
    v                  toggle between compact and spacious view

### Active Tasks List

    q/esc/ctrl+c       quit
    o/a                add task below cursor
    O                  add task above cursor
    I                  add task at the top
    A                  add task at the end
    u                  update task summary
    ⏎                  move task to the top
    J                  move task one position down
    K                  move task one position up

### Task Details Pane

    h/←/→/l            move backwards/forwards when in the task details view
    y                  copy current task's context to system clipboard
    B                  open all bookmarks added to current task

### Task Bookmarks List

    ⏎                  open URL in browser

### Quick Filter List

    ⏎                  pre-populate task list's search prompt with chosen prefix
