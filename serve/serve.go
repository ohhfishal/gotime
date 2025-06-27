package serve

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type EntryHandler interface {
	CreateEntry(entry.Entry) error
	GetAllEntries() ([]entry.Entry, error)
}

func Serve(ctx context.Context, logger *slog.Logger, handler EntryHandler, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		page().Render(ctx, w)
	})

	mux.HandleFunc("GET /api/v1/entry", func(w http.ResponseWriter, r *http.Request) {
		//   GET /api/v1/entry (Get's the HTML to render the Schedule)
		//   TODO: Implement
	})

	mux.HandleFunc("POST /api/v1/entry", func(w http.ResponseWriter, r *http.Request) {
		//   POST /api/v1/entry (Post's the entry)
		//   TODO: Implement
	})

	// TODO: Have this return a nice page instead!
	mux.Handle("/", http.NotFoundHandler())

	server := &http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("closing server",
				"error", err,
			)
		}
	}()

	logger.Info("starting server", "port", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
