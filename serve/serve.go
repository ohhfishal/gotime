package serve

import (
	"context"
	_ "embed"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ohhfishal/gotime/assets"
	"github.com/ohhfishal/gotime/entry"
)

//go:embed openapi.yaml
var spec []byte

type EntryHandler interface {
	CreateEntry(entry.Entry) error
	GetAllEntries(...entry.Option) ([]entry.Entry, error)
}

const errorHTML = `
<html>
  <h1>
    {{ .Error }}
  </h1>
</html>
`

var errorTemplate = template.Must(template.New("error").Parse(errorHTML))

func Serve(ctx context.Context, logger *slog.Logger, handler EntryHandler, port string) error {
	r := chi.NewRouter()

	r.Use(loggingMiddleware(logger))
	r.Use(middleware.Recoverer)

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(spec)
	})

	r.Mount("/openapi", http.StripPrefix("/openapi", swaggerui.Handler(spec)))
	r.Mount("/assets", http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))))

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets.Assets, "img/favicon.ico")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Have this be handled via HTMX and GET endpoint
		entries, err := handler.GetAllEntries(entry.Today())
		if err != nil {
			// TODO: Make a template for this
			w.WriteHeader(http.StatusInternalServerError)
			errorTemplate.Execute(w, map[string]string{"Error": err.Error()})
			return
		}
		MainPage(
			entry.Summarize(entries),
		).Render(r.Context(), w)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/entry", func(w http.ResponseWriter, r *http.Request) {
			entries, err := handler.GetAllEntries(entry.Today())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorTemplate.Execute(w, map[string]string{"Error": err.Error()})
				return
			}
			Details(
				entry.Summarize(entries),
			).Render(r.Context(), w)

			// TODO: Filter entries based on some parameters
			// TODO: Have this return HTML or JSON based on the header
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(entries)
		})

		r.Post("/entry", func(w http.ResponseWriter, r *http.Request) {
			var entryData entry.Entry
			if err := json.NewDecoder(r.Body).Decode(&entryData); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			// TODO: Actually validate the entry!

			if err := handler.CreateEntry(entryData); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			// TODO: Have this return HTML or JSON or text/plain based on the header
			w.WriteHeader(http.StatusOK)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
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

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			logger.Info("replied to request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.statusCode,
				"duration", time.Since(start).String(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
