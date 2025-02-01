package report

import (
	"fmt"
	"strings"
	"time"
)

var funcMap = map[string]any{
	"time":     Time,
	"duration": Duration,
}

func Time(t time.Time) (string, error) {
	return t.Format("15:04"), nil
}

func Duration(duration time.Duration) (string, error) {
	duration = duration.Round(time.Minute)
	if duration.Hours() > 24 {
		return ``, fmt.Errorf("duration: too big (24h): %s", duration)
	}
	if duration.Seconds() == 0 {
		return duration.String(), nil
	}
	return strings.TrimSuffix(duration.String(), "0s"), nil
}
