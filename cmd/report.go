package cmd

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
)

// TODO: Add some options!
// - [ ] Different time formatting? Turn on feature booleans?
// - [ ] Timezones??
// - [ ] Extra metadata?

// TODO: Expose a flag to use report.WithBroadCategoryBreakdown
//       Haven't done it since that requires more touching of the templates

type ReportCmd struct {
	// TODO: Times are tricky... Make sure this is the same timezone as time.Now()
	Start           time.Time `short:"s" optional:"" format:"2006-01-02" help:"Date to start report from (default: today)"`
	End             time.Time `short:"e" optional:"" format:"2006-01-02" help:"Date to end report from (default: time.Now)"`
	UseRootCategory bool      `help:"Enable to report final time worked the root of a category (example:'proj/task' would be reported as 'proj''" env:"GOTIME_USE_ROOT_CATEGORY"`
	DurationFormat  string    `enum:"default,hour" default:"default" help:"How to format duration in output (values: default,hour)" env:"GOTIME_DURATION_FORMAT"`
	// TODO: Make this a cleaner API?
	Output   report.OutputFormat `short:"o" enum:"default,markdown,html,json" default:"default" help:"Premade formats (default,markdown,html,json)"`
	Template string              `type:"existingfile" help:"File text/template to use in making report (Overrides --output)"`
	// TODO: This gets tricky since it could also be time...
	// TODO: Also probably should be an ENV?
	Until           time.Duration   `short:"u" help:"Desired duration of work to log"`
	templateContent string          `kong:"-"`
	reportOptions   []report.Option `kong:"-"`
}

func (cmd *ReportCmd) Run(stdout io.Writer, log string) error {
	entries, err := entry.ReadAllFromFile(log)
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	return report.Report(stdout,
		entries,
		cmd.templateContent,
		cmd.Start,
		cmd.End,
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

	// TODO: Maybe have this be conditional. It's now if start.IsZero otherwise default is start + 1d?
	if cmd.End.IsZero() {
		cmd.End = time.Now()
	}

	if cmd.Output == report.OutputFormatJSON {
		cmd.reportOptions = append(cmd.reportOptions, report.UseJSON)
	} else if content, err := cmd.Output.Template(); err == nil {
		cmd.templateContent = content
	}

	if cmd.Template != `` {
		bytes, err := os.ReadFile(cmd.Template)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
		cmd.templateContent = string(bytes)
	}

	// TODO: Use the new options
	cmd.reportOptions = append(cmd.reportOptions, report.WithDurationFormat(cmd.DurationFormat))
	if cmd.Until.Seconds() != float64(0) {
		cmd.reportOptions = append(cmd.reportOptions, report.WithUntil(cmd.Until))
	}
	if cmd.UseRootCategory {
		cmd.reportOptions = append(cmd.reportOptions, report.WithBroadCategoryBreakdown())
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
