

# Change Log

## [Upcoming] - 2025-XX-XX
### Added
- `tui` option `debug`

### Fixed
- Bug with `tui`'s report not updating after the first time
- Standard `report` to show the estimated end time when using `--until`

## [v0.7.0] - 2025-06-17
### **Breaking** ⚠️
- Moved `until` to be an option in `report`
- Renamed env `DURATION_FORMAT` to `GOTIME_DURATION_FORMAT`

### Added
- `report` option to treat category paths ("project/task") as their root ("project")
- `report` output `json`
- `tui` proof of concept

## [v0.6.0] - 2025-05-20
### **Breaking** ⚠️
- Changed log file to be a CSV. To transition, use the following command: `sed 's/: /, /g' $FILE`.
- Removed `append` and `continue` commands

### Added
- `report` html template
- `log` "-" category as an alternative to `append`
 
## [v0.5.0] - 2025-05-13
### Added
- `report` markdown template
- `report` duration formatting options

### Changed
- `log` to allow arbitrary log times
- `report` to use dates rather than duration offsets

## [v0.4.0] - 2025-02-28

### Added
- `until` command

### Changed
- `report` standard formatting spacing
- `log` to be the default command
- Flattened command tree

## [v0.3.1] - 2025-02-14

### Added
- Initial `v1` design document
- Help messages to all commands
- `README.md` improvements
- Total logged time to standard `report` template
- Special "out" log category

### Changed
- `$GOTIME_REPORT_TEMPLATE` to be a flag in `report`
- Default for `$GOTIME_LOG_PATH` 

## [v0.3.0] - 2025-02-08

### Added
- Help messages for `report` and `log` commands
- New `log` subcommands:
  - `append`
  - `continue`
  - `new` (default)

## [v0.2.0] - 2025-02-02

### Added
- Day filtering
- Report using `text/template`
 
## [v0.1.0] - 2025-01-31
 
- Initial release
 
[upcoming]: https://github.com/ohhfishal/gotime/compare/v0.7.0...HEAD
[v0.7.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.7.0
[v0.6.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.6.0
[v0.5.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.5.0
[v0.4.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.4.0
[v0.3.1]: https://github.com/ohhfishal/gotime/releases/tag/v0.3.1
[v0.3.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.3.0
[v0.2.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.2.0
[v0.1.0]: https://github.com/ohhfishal/gotime/releases/tag/v0.1.0
