# Design Document

## Purpose
Provide a useful CLI tool for managing ones daily work journal.

## Terms
- `entry`: A journal entry that describes what was done between a range of time (Ex: I read from 1-2pm).

# V1 Goals
## Commands
- [X] `log`: Create a new entry in your daily log
- [X] `append`: Log another entry using your current entry.
- [X] `resume`: Continue using a given category.
- [X] `report`: Prints out a human-readable report.

## Improvements
- [ ] Better docs
- [ ] Improved `report` functionality 
  - [ ] More templates (Markdown)
  - [ ] Condense `until` command to be part of `report`)
  - [ ] Subcategories `project/task`


# Futher Goals (Not required for V1)
- [ ] `export`: Exports data to a standard format (Ex: `json`).
