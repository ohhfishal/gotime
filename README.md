# gotime
Go CLI to track time spent on tasks


## Basic Examples
```
go build
./gotime log ProjectA Here is a description
./gotime report
```

## ENV Configuration

### Example
```bash
export GOTIME_LOG="/home/user/.gotime/schedule.log"
```

| Name | Description | Default | 
| ---- | ----------- | ------- | 
| `GOTIME_LOG` | Path to the file to read/write tasks in (BAD DEFAULT SET THIS) | `gotime.log` |
| `GOTIME_NOFILE` | When "stdout" `log` prints out to stdout rather than a file (for debug/tests) | "" | 
| `GOTIME_REPORT_TEMPLATE` | Override for `report` see `report/templates/standard.tpl` | | 
