package report

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

type UntilConfig struct {
	Current time.Duration
	Total   time.Duration
	Left    time.Duration
}

func ReportUntil(stdout io.Writer, templatePath string, config UntilConfig) error {
	config.Left = config.Total - config.Current
	duration, err := Duration(config.Left)
	if err != nil {
		return fmt.Errorf(`parsing until duration: %w`, err)
	}

	if templatePath == `` {
		fmt.Fprintln(stdout, duration)
		return nil
	}

	bytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("reading template file: %w", err)
	}

	tmpl, err := template.New("report-until-template").Funcs(defaultFuncMap).Parse(string(bytes))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = tmpl.Execute(stdout, config)
	if err != nil {
		return fmt.Errorf("printing template: %w", err)
	}
	return nil
}
