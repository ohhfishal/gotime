package cmd

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type ReportCmd struct {
	Back    time.Duration `default:"0h0m" help:"Offset to go back from 00:00 today (today@00:00 - offset)."`
	Forward time.Duration `default:"24h" help:"Offset to go forward from 00:00 today (today@00:00 + offset)."`
}

func (cmd *ReportCmd) Run(cfg Config) error {
	path := cfg.LogPath()
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`can not open "%s": %w`, path, err)
	}

	// Resolve when to start and end
	now := time.Now().Format(time.DateOnly)
	today, err := time.Parse(time.DateOnly, now)
	if err != nil {
		return fmt.Errorf("parsing today: %w", err)
	}
	start := today.Add(-cmd.Back)
	end := today.Add(cmd.Forward)

	// TODO: Add in config to customize the file format read
	entries, err := entry.ReadAll(file)
	if err != nil {
		return fmt.Errorf(`parsing file "%s": %w`, path, err)
	}
	filtered := entry.Filter(entries, start, end)
	slices.SortFunc(filtered, entry.Compare)
	return report.Report(report.Config{
		Stdout:   cfg.Stdout,
		Template: cfg.Getenv("GOTIME_REPORT_TEMPLATE"),
	}, filtered)
}

func Report(cfg Config, args ...string) error {
	var cmd ReportCmd
	return RunCmd(cfg, &cmd, args...)
}
