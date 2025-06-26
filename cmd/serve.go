package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ohhfishal/gotime/report"
)

type ServeCmd struct {
	Port   string       `short:"p" default:":8080" help:"Port to run server on"`
	Report ReportCmd    `embed:""`
	log    string       `kong:"-"`
	logger *slog.Logger `kong:"-"`
}

func (cmd *ServeCmd) AfterApply() error {
	if err := cmd.Report.AfterApply(); err != nil {
		return err
	}
	cmd.Report.Output = report.OutputFormatHTML
	return nil
}

func (cmd *ServeCmd) Run(stdout io.Writer, log string) error {
	cmd.logger = slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // TODO: Make configurable
	}))
	cmd.log = log

	// TODO: Don't have this all use global http
	mux := http.NewServeMux()
	mux.Handle("/", cmd)
	cmd.logger.Info("starting server", "port", cmd.Port)

	// TODO: Respect CTRL + C
	return http.ListenAndServe(cmd.Port, mux)
}

func (handler *ServeCmd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("responding to request")
	w.Header().Set("Content-Type", "text/html")
	html, err := handler.getReport()
	if err != nil {
		// TODO: Resolve this
		panic(err)
	}
	fmt.Fprint(w, html)
}

func (handler *ServeCmd) getReport() (string, error) {
	handler.Report.End = time.Now()
	var stdout strings.Builder
	if err := handler.Report.Run(&stdout, handler.log); err != nil {
		// TODO: Resolve this better
		return ``, err
	}
	return stdout.String(), nil
}
