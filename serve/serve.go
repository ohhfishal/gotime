package serve

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/ohhfishal/gotime/entry"
)

//go:embed openapi.yaml
var spec []byte

type EntryHandler interface {
	CreateEntry(entry.Entry) error
	GetAllEntries() ([]entry.Entry, error)
}

func Serve(ctx context.Context, logger *slog.Logger, handler EntryHandler, port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Write(spec)
	})
	mux.Handle("GET /openapi/", http.StripPrefix("/openapi", swaggerui.Handler(spec)))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		page().Render(ctx, w)
	})

	mux.HandleFunc("GET /api/v1/entry", func(w http.ResponseWriter, r *http.Request) {
		//   GET /api/v1/entry (Get's the HTML to render the Schedule)
		//   TODO: Implement
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, http.StatusText(http.StatusNotImplemented))
	})

	mux.HandleFunc("POST /api/v1/entry", func(w http.ResponseWriter, r *http.Request) {
		//   POST /api/v1/entry (Post's the entry)
		//   TODO: Implement
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, http.StatusText(http.StatusNotImplemented))
	})

	// TODO: Have this return a nice page instead!
	mux.Handle("/", http.NotFoundHandler())

	var responseHandler http.Handler
	responseHandler = NewLoggingMiddleware(logger, mux)

	server := &http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      responseHandler,
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

func NewLoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		// TODO: Make better
		logger.Info("replied to request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start).String(),
		)
	})
}

// TODO: Should write the message in here as well
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
