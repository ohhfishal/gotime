package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type UntilCmd struct {
	Duration time.Duration `arg:"" required:""`
	Exclude  []string      `default:"out" help:"Categories to exclude from total"`
	Category string        `default:"total" help:"Category to check against duration"`
	Template string        `type:"existingfile" placeholder:"PATH" help:"Override output using a text/template"`
}

func (cmd *UntilCmd) Run(stdout io.Writer, log string, nowFunc func() time.Time) error {
	entries, err := entry.ReadAllFromFile(log)
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	now, err := Standardize(nowFunc())
	if err != nil {
		return fmt.Errorf("parsing now: %w", err)
	}

	today, err := Today(now)
	if err != nil {
		return fmt.Errorf("parsing today: %w", err)
	}

	filtered := entry.Filter(entries, today, now)
	durations := entry.DurationMap(filtered, now)
	for _, exclude := range cmd.Exclude {
		delete(durations, exclude)
	}

	durations[`total`] = entry.Total(durations)
	total, ok := durations[cmd.Category]
	if !ok {
		return fmt.Errorf(`unknown category: ""%s`, cmd.Category)
	}

	if err := report.ReportUntil(stdout, cmd.Template, report.UntilConfig{
		Current: total,
		Total:   cmd.Duration,
	}); err != nil {
		return fmt.Errorf(`printing report: %w`, err)
	}
	return nil
}

func Standardize(timestamp time.Time) (time.Time, error) {
	now := timestamp.Format(time.DateTime)
	today, err := time.Parse(time.DateTime, now)
	if err != nil {
		return time.Time{}, err
	}
	return today, nil
}

// TODO: May be shared by other code
func Today(timestamp time.Time) (time.Time, error) {
	now := timestamp.Format(time.DateOnly)
	today, err := time.Parse(time.DateOnly, now)
	if err != nil {
		return time.Time{}, err
	}
	return today, nil
}
