# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Allow duplicating tasks

## [v0.6.0] - May 02, 2025

### Changes

- Use full width for context and task details
- Increase upper threshold for context length to 1 MB (was 4KB before)
- Preserve newlines in rendered context

## [v0.5.1] - Aug 13, 2024

### Fixed

- Fixed issue where omm would panic when tasks summaries of certain lengths were
  entered (on certain platforms)

## [v0.5.0] - Aug 03, 2024

### Added

- URIs with custom schemes are considered as task bookmarks. eg.
    - `spotify:track:4fVBFyglBhMf0erfF7pBJp`
    - `obsidian://open?vault=VAULT&file=FILE`
- Circular navigation for lists

## [v0.4.3] - Jul 28, 2024

### Added

- Flag for changing deletion behavior

### Changes

- omm asks for confirmation before deleting a task by default

## [v0.4.2] - Jul 26, 2024

### Fixed

- Fixed issue where pager didn't respond to arrow keys

## [v0.4.0] - Jul 26, 2024

### Added

- Markdown in task context is rendered with syntax highlighting
- Task Lists in "compact" mode highlight prefixes
- Tasks Lists can be filtered
- Quick filters based on task prefixes can be applied
- Prefix can be chosen during task creation/update 
- Keymap to move an active task to the end of the list
- "updates" subcommand

### Changes

- Task Lists in "compact" mode show more than 9 tasks at a time, when maximised

### Removed

- Keymaps ([2-9]) for the Active Tasks list to move task at a specific index to
  the top

## [v0.3.1] - Jul 19, 2024

### Changes

- URLs in a task's summary are also considered as bookmarks

## [v0.3.0] - Jul 19, 2024

### Added

- The ability to quickly open URLs present in a task's context
- Support for providing configuration via a TOML file
- The ability to copy a task's context

## [v0.2.2] - Jul 15, 2024

## Fixed

- Fixed issue where closing the "Task Details" pane would move the active task
  list to the next page

## [v0.2.1] - Jul 15, 2024

## Fixed

- Fixed issue where omm's database would be stored in an incorrect location on
  Windows

## [v0.2.0] - Jul 14, 2024

### Added

- Added "task context", which can be used for additional details for a task that
  don't fit in the summary
- A new list density mode
- "Task Details" pane
- An onboarding guide

### Changed

- Task lists now highlight prefixes in task summaries, when provided

## [v0.1.0] - Jul 09, 2024

### Added

- Initial release

[unreleased]: https://github.com/dhth/omm/compare/v0.6.0...HEAD
[v0.6.0]: https://github.com/dhth/omm/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/dhth/omm/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/dhth/omm/compare/v0.4.3...v0.5.0
[v0.4.3]: https://github.com/dhth/omm/compare/v0.4.2...v0.4.3
[v0.4.2]: https://github.com/dhth/omm/compare/v0.4.0...v0.4.2
[v0.4.0]: https://github.com/dhth/omm/compare/v0.3.1...v0.4.0
[v0.3.1]: https://github.com/dhth/omm/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/dhth/omm/compare/v0.2.2...v0.3.0
[v0.2.2]: https://github.com/dhth/omm/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/dhth/omm/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/dhth/omm/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/dhth/omm/commits/v0.1.0/
