package report

import (
	"fmt"
	"strings"
	"time"
)

// TODO: Add tests for these
var defaultFuncMap = map[string]any{
	"time":       Time,
	"date":       Date,
	"duration":   Duration,
	"truncRight": TruncRightWith,
}

func TruncRightWith(size int, suffix, msg string) (string, error) {
	if len(suffix) > size {
		return ``, fmt.Errorf(`suffix "%s" is too large (%d)`, suffix, size)
	}
	if len(msg) <= size {
		return msg, nil
	}
	return msg[:size-len(suffix)] + suffix, nil
}

func Date(t time.Time) (string, error) {
	return t.Format("2006-01-02"), nil
}

func Time(t time.Time) (string, error) {
	return t.Format("15:04"), nil
}

func Duration(duration time.Duration) (string, error) {
	duration = duration.Round(time.Minute)
	if duration.Seconds() == 0 {
		return duration.String(), nil
	}
	return strings.TrimSuffix(duration.String(), "0s"), nil
}

func DurationHour(t time.Duration) (string, error) {
	return fmt.Sprintf(`%.1f`, t.Hours()), nil
}
