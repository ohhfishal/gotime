package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/report"
	"github.com/ohhfishal/gotime/serve"
)

type ServeCmd struct {
	Port   string    `short:"p" default:"8080" help:"Port to run server on"`
	Report ReportCmd `embed:""`
}

func (cmd *ServeCmd) AfterApply() error {
	if err := cmd.Report.AfterApply(); err != nil {
		return err
	}
	cmd.Report.Output = report.OutputFormatHTML
	return nil
}

func (cmd *ServeCmd) Run(stdout io.Writer, log string) error {
	logger := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // TODO: Make configurable
	}))

	entryHandler, err := entry.NewFileHandler(log)
	if err != nil {
		return fmt.Errorf(`connecting to file handler: %w`, err)
	}

	// TODO: Actually use a handler for getting/retriving entries
	return serve.Serve(context.TODO(), logger, entryHandler, cmd.Port)
}
