package cmd

import (
	"fmt"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type UntilCmd struct {
	Duration time.Duration `arg: "" required: ""`
	Exclude  []string      `default:"out" help:"Categories to exclude from total"`
	Category string        `default:"total" help:"Category to check against duration"`
	Template string        `type:"path" help:"Override output using a text/template"`
}

func (cmd *UntilCmd) Run(cfg Config) error {
	entries, err := cfg.GetAllEntries()
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	today, err := cfg.Today()
	if err != nil {
		return fmt.Errorf("parsing today: %w", err)
	}

	now, err := cfg.Now()
	if err != nil {
		return fmt.Errorf("parsing current time: %w", err)
	}

	filtered := entry.Filter(entries, *today, *now)
	durations := entry.DurationMap(filtered, *now)
	for _, exclude := range cmd.Exclude {
		delete(durations, exclude)
	}

	durations[`total`] = entry.Total(durations)
	total, ok := durations[cmd.Category]
	if !ok {
		return fmt.Errorf(`unknown category: ""%s`, cmd.Category)
	}

  if err := report.ReportUntil(cfg.Stdout, cmd.Template, report.UntilConfig{
    Current: total, 
    Total: cmd.Duration,
  }); err != nil {
    return fmt.Errorf(`printing report: %w`, err)
  }
	return nil
}

func Until(cfg Config, args ...string) error {
	var cmd UntilCmd
	return RunCmd(cfg, cmd, args...)
}
