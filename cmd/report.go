package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type ReportCmd struct {
	Back    time.Duration `default:"0h0m" help:"Offset to go back from 00:00 today (today@00:00 - offset)."`
	Forward time.Duration `default:"24h" help:"Offset to go forward from 00:00 today (today@00:00 + offset)."`
  Template string `optional:"" type:"existingfile" help:"text/templat to use in making report"`
}

func (cmd *ReportCmd) Run(cfg Config) error {
	// Resolve when to start and end
	now := time.Now().Format(time.DateOnly)
	today, err := time.Parse(time.DateOnly, now)
	if err != nil {
		return fmt.Errorf("parsing today: %w", err)
	}
	start := today.Add(-cmd.Back)
	end := today.Add(cmd.Forward)

	entries, err := cfg.GetAllEntries()
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}
	filtered := entry.Filter(entries, start, end)
	slices.SortFunc(filtered, entry.Compare)
	return report.Report(report.Config{
		Stdout:   cfg.Stdout,
		Template: cmd.Template,
	}, filtered)
}

func Report(cfg Config, args ...string) error {
	var cmd ReportCmd
	return RunCmd(cfg, &cmd, args...)
}
