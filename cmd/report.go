package cmd

import (
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type ReportCmd struct {
	Back     time.Duration `default:"0h0m" help:"Offset to go back from 00:00 today (today@00:00 - offset)."`
	Forward  time.Duration `default:"24h" help:"Offset to go forward from 00:00 today (today@00:00 + offset)."`
	Template string        `optional:"" type:"existingfile" help:"text/template to use in making report"`
}

func (cmd *ReportCmd) Run(stdout io.Writer, entrySet EntrySet, now func() time.Time) error {
	// Resolve when to start and end
	today, err := Today(now())
	if err != nil {
		return fmt.Errorf("parsing today: %w", err)
	}

	start := today.Add(-cmd.Back)
	end := today.Add(cmd.Forward)

	entries, err := entrySet.GetAll()
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	filtered := entry.Filter(entries, start, end)
	slices.SortFunc(filtered, entry.Compare)
	return report.Report(report.Config{
		Stdout:   stdout,
		Template: cmd.Template,
	}, filtered)
}
