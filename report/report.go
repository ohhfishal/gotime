package report

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
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
}

func (config Config) getTemplate() (string, error) {
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

func annotate(entries []entry.Entry) (*Metadata, error) {
	schedule := []Entry{}
	totals := map[string]time.Duration{}
	for i, cur := range entries {
		next := entry.Entry{Time: time.Now()}
		if i < (len(entries) - 1) {
			next = entries[i+1]
		}

		newEntry, err := annotateEntry(cur, next)
		if err != nil {
			return nil, err
		}
		schedule = append(schedule, newEntry)
		if _, ok := totals[newEntry.Category]; ok {
			totals[newEntry.Category] += newEntry.Duration
		} else {
			totals[newEntry.Category] = newEntry.Duration
		}
	}
	return &Metadata{
		Schedule:   schedule,
		Categories: totals,
	}, nil
}

func annotateEntry(entry, next entry.Entry) (Entry, error) {
	newEntry := Entry{
		Category: entry.Category,
		Note:     entry.Note,
		Time:     entry.Time,
		Duration: next.Time.Sub(entry.Time),
	}
	// TODO: Handle negative time case
	return newEntry, nil
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

	tmpl, err := template.New("report-template").Funcs(sprig.FuncMap()).Funcs(funcMap).Parse(templateString)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = tmpl.Execute(config.Stdout, entries)
	if err != nil {
		return fmt.Errorf("printing template: %w", err)
	}
	return nil
}
