package report

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"
)

// TODO: Add tests for these
var defaultFuncMap = map[string]any{
	"time":       Time,
	"date":       Date,
	"duration":   Duration,
	"truncRight": TruncRightWith,
}

//go:embed templates/standard.tpl
var defaultTemplate string

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
		return defaultTemplate, nil
	default:
		return ``, errors.New(`invalid format`)
	}
}

func TruncRightWith(size int, suffix, msg string) (string, error) {
	if len(suffix) > size {
		return ``, fmt.Errorf(`suffix "%s" is too large (%d)`, suffix, size)
	}
	if len(msg) <= size {
		return msg, nil
	}
	return msg[:size-len(suffix)] + suffix, nil
}

func Date(t time.Time) (string, error) {
	return t.Format("2006-01-02"), nil
}

func Time(t time.Time) (string, error) {
	return t.Format("15:04"), nil
}

func Duration(duration time.Duration) (string, error) {
	duration = duration.Round(time.Minute)
	if duration.Seconds() == 0 {
		return duration.String(), nil
	}
	return strings.TrimSuffix(duration.String(), "0s"), nil
}

func DurationHour(t time.Duration) (string, error) {
	return fmt.Sprintf(`%.1f`, t.Hours()), nil
}
