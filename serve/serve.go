package serve

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ohhfishal/gotime/entry"
)

type EntryHandler interface {
	CreateEntry(entry.Entry) error
	GetAllEntries() ([]entry.Entry, error)
}

func Serve(ctx context.Context, logger *slog.Logger, handler EntryHandler) error {
	mux := http.NewServeMux()
	mux.Handle("/", cmd)
	cmd.logger.Info("starting server", "port", cmd.Port)

	// TODO: Respect the context!
	return http.ListenAndServe(cmd.Port, mux)
}

// TODO: Implement the rest
// Routes:
//   POST /api/v1/entry (Post's the entry)
//   GET /api/v1/entry (Get's the HTML to render the Schedule)
//   GET / (Renders main page)
//
