# gotime
[![Go Report Card](https://goreportcard.com/badge/github.com/ohhfishal/gotime)](https://goreportcard.com/report/github.com/ohhfishal/gotime)

Go CLI to track time spent on tasks

## Install

```bash
go install github.com/ohhfishal/gotime@latest
```

## Build
```bash
git clone git@github.com:ohhfishal/gotime.git
go build
```

## Basic Examples
If the following were run throughout the day:
```
./gotime log project Working on Bug
./gotime log out/lunch
./gotime continue project
./gotime append Meeting
./gotime log out Done for the day
./gotime report
```

The final `gotime report` may produce:
```
 Category       Start   Length      Note
 --------       -----   ------      ----
 project        08:00   3h30m       Working on Bug
 out/lunch      11:30   30m
 project        12:00   3h0m        Cont: Working on Bug
 project        15:00   1h30m       Meeting
 out            16:30   0s          Done for the day


 Category       Length
 --------       ------
 out            0s
 out/lunch      30m
 project        8h0m

```

## Special categories

`out` is treated as a special entry category. It's time is not used in calculating the `total` field in the standard report.

## ENV Configuration


| Name | Description | Default |
| ---- | ----------- | ------- |
| `GOTIME_LOG` | Filepath to log file used to `log` and `report` tasks. | `~/.config/gotime.csv` |
| `DURATION_FORMAT` | Enum that controls the way durations are formatted (default or hour) | `default` |

### Example
```bash
export GOTIME_LOG="~/.config/gotime.csv"
export DURATION_FORMAT="hour"
```

