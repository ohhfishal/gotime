package report

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

//go:embed templates/standard.tpl
var defaultTemplate string

var defaultTemplates = map[string]string{
	"": defaultTemplate,
}

type Config struct {
	Stdout   io.Writer
	Template string
}

type Entry struct {
	Time     time.Time
	Category string
	Note     string
	Duration time.Duration
}

// TODO: Use this instead of the []Entry
type Metadata struct {
	Schedule   []Entry
	Categories map[string]time.Duration
	Total      time.Duration
}

func (config Config) getTemplate() (string, error) {
	// TODO: Remove this and just read the path from the args
	// Check for default templates
	tmpl, ok := defaultTemplates[config.Template]
	if !ok {
		return config.loadTemplate(config.Template)
	}
	return tmpl, nil
}

func (config Config) loadTemplate(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return ``, fmt.Errorf("reading file: %w", err)
	}
	return string(bytes), nil
}

func deviceNow() time.Time {
	final, err := time.Parse(time.DateTime, time.Now().Format(time.DateTime))
	if err != nil {
		// TODO: Consider if we keep this as a panic
		panic(fmt.Errorf("could not get device time: %w", err))
	}
	return final

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
	}, nil
}

func ReportUntil(stdout io.Writer, template string, current, total time.Duration) error {
	duration, err := Duration(total - current)
	if err != nil {
		return fmt.Errorf(`parsing until duration: %w`, err)
	}

	if template != `` {
		// TODO: Handle the template here
		return fmt.Errorf(`text/template not implemented for "log until"`)
	}
	fmt.Fprintln(stdout, duration)
	return nil
}

func Report(config Config, originals []entry.Entry) error {
	entries, err := annotate(originals)
	if err != nil {
		return fmt.Errorf("formatting entries: %w", err)
	}
	templateString, err := config.getTemplate()
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}

	tmpl, err := template.New("report-template").Funcs(funcMap).Parse(templateString)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = tmpl.Execute(config.Stdout, entries)
	if err != nil {
		return fmt.Errorf("printing template: %w", err)
	}
	return nil
}
