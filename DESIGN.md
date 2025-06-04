# Design Document

## Purpose
Provide a useful CLI tool for managing ones daily work journal.

## Terms
- `entry`: A journal entry that describes what was done between a range of time (Ex: I read from 1-2pm).

# V1 Goals
## Commands
- [X] `log`: Create a new entry in your daily log
- [X] `report`: Prints out a human-readable report.
- [ ] `tui`: Run gotime as a TUI.

## Improvements
- [ ] Better docs
  - [ ] https://pkg.go.dev/github.com/ohhfishal/gotime
  - [ ] Video example of `tui`
  - [ ] Redesigned `README.md`
- [ ] Improved `report` functionality 
  - [X] More templates (Markdown)
  - [X] Condense `until` command to be part of `report`)
  - [ ] Subcategories `project/task`

# Futher Goals (Not required for V1)
- [X] `export`: Exports data to a standard format (Ex: `json`). (**NOTE:** Counting moving to a CSV based log file as a means of a [credible exit](https://newsletter.squishy.computer/p/credible-exit))
