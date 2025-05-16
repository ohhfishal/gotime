package report

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

//go:embed templates/standard.tpl
var defaultTemplate string

// TODO: Make this template look good as plain text
//
//go:embed templates/markdown.tpl
var markdownTemplate string

//go:embed templates/html.tpl
var htmlTemplate string

type OutputFormat string

const (
	OutputFormatDefault  = ""
	OutputFormatMarkdown = "markdown"
	OutputFormatHTML     = "html"
)

var durationFormats map[string]any = map[string]any{
	"default": Duration,
	"hour":    DurationHour,
}

func (format OutputFormat) Template() (string, error) {
	switch format {
	case OutputFormatHTML:
		return htmlTemplate, nil
	case OutputFormatMarkdown:
		return markdownTemplate, nil
	case OutputFormatDefault:
		fallthrough
	default:
		return ``, errors.New(`invalid format`)
	}
}

var defaultTemplates = map[string]string{
	"": defaultTemplate,
}

// TODO: Add some options!
// - [ ] Different time formatting? Turn on feature booleans?
// - [ ] Extra metadata?
// - [ ] Move the templateString here? WithTemplate(string)
type ReportOption func(*Metadata) error

func WithDurationFormat(format string) func(*Metadata) error {
	return func(metadata *Metadata) error {
		f, ok := durationFormats[format]
		if !ok {
			return fmt.Errorf(`invalid duration format: %s`, format)
		}
		metadata.funcs["duration"] = f
		return nil
	}
}

type Metadata struct {
	Schedule   []Entry
	Categories map[string]time.Duration
	Total      time.Duration
	funcs      map[string]any
	Time       time.Time
}

type Entry struct {
	Time     time.Time
	Category string
	Note     string
	Duration time.Duration
}

func Report(
	stdout io.Writer,
	templateString string,
	start time.Time,
	originals []entry.Entry,
	options ...ReportOption,
) error {
	metadata, err := annotate(originals)
	if err != nil {
		return fmt.Errorf("formatting metadata: %w", err)
	}
	metadata.Time = start

	for i, opt := range options {
		if err := opt(metadata); err != nil {
			return fmt.Errorf(`applying option %d: %w`, i, err)
		}
	}

	if templateString == `` {
		templateString = defaultTemplate
	}

	tmpl, err := template.
		New("report-template").
		Funcs(metadata.funcs).
		Parse(templateString)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = tmpl.Execute(stdout, metadata)
	if err != nil {
		return fmt.Errorf("printing template: %w", err)
	}
	return nil
}

func annotate(entries []entry.Entry) (*Metadata, error) {
	schedule := []Entry{}
	totals := map[string]time.Duration{}
	var total time.Duration
	for i, cur := range entries {
		newEntry := Entry{
			Category: cur.Category,
			Note:     cur.Note,
			Time:     cur.Time,
		}
		if i < (len(entries) - 1) {
			newEntry.Duration = entries[i+1].Time.Sub(cur.Time)
		} else {
			newEntry.Duration = deviceNow().Sub(cur.Time)
		}

		schedule = append(schedule, newEntry)
		if _, ok := totals[newEntry.Category]; ok {
			totals[newEntry.Category] += newEntry.Duration
		} else {
			totals[newEntry.Category] = newEntry.Duration
		}

		// TODO: Make this logic more generalizable
		if newEntry.Category != `out` {
			total += newEntry.Duration
		}
	}
	return &Metadata{
		Categories: totals,
		Schedule:   schedule,
		Total:      total,
		funcs:      defaultFuncMap,
	}, nil
}

// TODO: Should really remove this function
func deviceNow() time.Time {
	final, err := time.Parse(time.DateTime, time.Now().Format(time.DateTime))
	if err != nil {
		// TODO: Consider if we keep this as a panic
		panic(fmt.Errorf("could not get device time: %w", err))
	}
	return final
}
