package cmd

import (
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

type ReportCmd struct {
	Start          time.Time `short:"s" optional:"" format:"2006-01-02" help:"Date to start report from (default: today)"`
	End            time.Time `short:"t" optional:"" format:"2006-01-02" help:"Date to end report from (default: start + 1d)"`
	DurationFormat string    `enum:"default,hour" default:"default" help:"How to format duration in output (values: default,hour)" env:"DURATION_FORMAT"`
	// TODO: Make this a cleaner API
	Output   report.OutputFormat `short:"o" enum:"default,markdown,html" default:"default" help:"Premade formats (default,markdown,html)"`
	Template string              `type:"existingfile" help:"File text/template to use in making report (Overrides --output)"`
	// TODO: This gets tricky since it could also be time...
	// TODO: Also probably should be an ENV?
	Until           time.Duration         `short:"u" help:"Desired duration of work to log"`
	templateContent string                `kong:"-"`
	reportOptions   []report.ReportOption `kong:"-"`
}

func (cmd *ReportCmd) Run(stdout io.Writer, log string) error {
	entries, err := entry.ReadAllFromFile(log)
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	filtered := entry.Filter(entries, cmd.Start, cmd.End)
	slices.SortFunc(filtered, entry.Compare)
	return report.Report(stdout,
		cmd.templateContent,
		cmd.Start,
		filtered,
		cmd.reportOptions...,
	)
}

func (cmd *ReportCmd) AfterApply() error {
	var err error
	if cmd.Start.IsZero() {
		cmd.Start, err = Today(time.Now())
		if err != nil {
			return fmt.Errorf(`calculating today's date: %w`, err)
		}
	}

	if cmd.End.IsZero() {
		cmd.End = cmd.Start.Add(24 * time.Hour)
	}

	if content, err := cmd.Output.Template(); err == nil {
		cmd.templateContent = content
	}

	if cmd.Template != `` {
		bytes, err := os.ReadFile(cmd.Template)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
		cmd.templateContent = string(bytes)
	}

	cmd.reportOptions = append(cmd.reportOptions, report.WithDurationFormat(cmd.DurationFormat))
	if cmd.Until.Seconds() != float64(0) {
		cmd.reportOptions = append(cmd.reportOptions, report.WithDesiredDurationLogged(cmd.Until))
	}

	return nil
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
