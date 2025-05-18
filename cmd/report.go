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
	Output          report.OutputFormat `short:"o" enum:"default,markdown,html" default:"default" help:"Premade formats (default,markdown,html)"`
	Template        string              `type:"existingfile" help:"File text/template to use in making report (Overrides --output)"`
	templateContent string              `kong:"-"`
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
		report.WithDurationFormat(cmd.DurationFormat),
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
	return nil
}
