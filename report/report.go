package report

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

func Report(
	stdout io.Writer,
	entries []entry.Entry,
	templateString string, start time.Time,
	end time.Time,
	options ...Option,
) error {

	args, err := annotate(entries, start, end, options...)
	if err != nil {
		return fmt.Errorf(`creating entries metadata: %w`, err)
	}

	if err := args.Print(stdout, templateString); err != nil {
		return fmt.Errorf(`printing: %w`, err)
	}
	return nil
}

type ScheduleTuple struct {
	Entry    entry.Entry
	Duration time.Duration
}

type ReportArgs struct {
	Schedule          []ScheduleTuple
	CategoryBreakdown map[string]time.Duration
	Total             time.Duration
	StartTime         time.Time
	EndTime           time.Time
	// Optional fields
	Until                  time.Duration
	templateFunctions      map[string]any
}

func annotate(
	entries []entry.Entry,
	start time.Time,
	end time.Time,
	opts ...Option,
) (*ReportArgs, error) {

	// Filter to entries during the time range and allow us to assume i+1 exists
	filteredEntries := append(
		entry.Filter(entries, start, end),
		entry.Entry{Time: end},
	)

	args := &ReportArgs{
		StartTime:              start,
		EndTime:                end,
		CategoryBreakdown:      map[string]time.Duration{},
		templateFunctions:      defaultFuncMap,
	}

	for i, entry := range filteredEntries[:len(filteredEntries)-1] {
		tuple := ScheduleTuple{
			Entry:    entry,
			Duration: filteredEntries[i+1].Time.Sub(entry.Time),
		}
		if _, ok := args.CategoryBreakdown[entry.Category]; ok {
			args.CategoryBreakdown[entry.Category] += tuple.Duration
		} else {
			args.CategoryBreakdown[entry.Category] = tuple.Duration
		}
		args.Schedule = append(args.Schedule, tuple)
		args.Total += tuple.Duration
	}

	for _, opt := range opts {
		if err := opt(args); err != nil {
			return nil, fmt.Errorf(`applying option: %w`, err)
		}
	}
	return args, nil
}

func (args ReportArgs) Print(stdout io.Writer, rawTemplate string) error {
	if rawTemplate == `` {
		rawTemplate = defaultTemplate
	}

	tmpl, err := template.
		New("report-template").
		Funcs(args.templateFunctions).
		Parse(rawTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = tmpl.Execute(stdout, args)
	if err != nil {
		return fmt.Errorf("printing template: %w", err)
	}
	return nil
}

// NOTE: Options that error should include their name in the error message
type Option func(*ReportArgs) error

func WithBroadCategoryBreakdown() Option {
	// NOTE: Doesn't have to be a higher-order function
	//       Only doing it for consistency with name/implementation
	return func(args *ReportArgs) error {
		roots := map[string]time.Duration{}
		for category, duration := range args.CategoryBreakdown {
			path := strings.Split(category, "/")
			if len(path) == 0 {
				return fmt.Errorf(`broad category breakdown: invalid category: %s`, category)
			}
			root := path[0]

			if _, ok := roots[root]; ok {
				roots[root] += duration
			} else {
				roots[root] = duration
			}
		}
		args.CategoryBreakdown = roots
		return nil
	}
}

func WithUntil(expected time.Duration) Option {
	return func(args *ReportArgs) error {
		args.Until = expected - args.Total
		return nil
	}
}

func WithDurationFormat(format string) Option {
	return func(args *ReportArgs) error {
		f, ok := durationFormats[format]
		if !ok {
			return fmt.Errorf(`duration format: invalid duration format: %s`, format)
		}
		args.templateFunctions["duration"] = f
		return nil
	}
}
