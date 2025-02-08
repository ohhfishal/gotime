# gotime
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
./gotime log continue project
./gotime log append Meeting
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


## ENV Configuration

### Example
```bash
export GOTIME_LOG="/home/user/.gotime/schedule.log"
```

| Name | Description | Default | 
| ---- | ----------- | ------- | 
| `GOTIME_LOG` | Path to the file to read/write tasks in (BAD DEFAULT SET THIS) | `gotime.log` |
| `GOTIME_REPORT_TEMPLATE` | Override for `report` see `report/templates/standard.tpl` | | 
