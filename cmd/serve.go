package cmd

import (
	"context"
	"io"
	"log/slog"

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
	// TODO: Actually use a handler for getting/retriving entries
	return serve.Serve(context.TODO(), logger, nil /* TODO: Implement*/, cmd.Port)
}

// func (handler *ServeCmd) getReport() (string, error) {
// 	handler.Report.End = time.Now()
// 	var stdout strings.Builder
// 	if err := handler.Report.Run(&stdout, handler.log); err != nil {
// 		// TODO: Resolve this better
// 		return ``, err
// 	}
// 	return stdout.String(), nil
// }
