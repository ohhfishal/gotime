package serve

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/ohhfishal/gotime/assets"
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

	mux.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		page().Render(ctx, w)
	})

	mux.Handle("GET /api/v1/entry", CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		entries, err := handler.GetAllEntries()
		if err != nil {
			return StatusWithError(http.StatusInternalServerError, err)
		}
		// TODO: Filter entries based on some parameters
		// TODO: Have this return HTML or JSON based on the header
		return JSON(http.StatusOK, entries)
	}))

	mux.Handle("POST /api/v1/entry", CustomHandler(func(w http.ResponseWriter, r *http.Request) http.Handler {
		//   POST /api/v1/entry (Post's the entry)
		//   TODO: Implement
		entry, err := decode[entry.Entry](r)
		if err != nil {
			return StatusWithError(http.StatusBadRequest, err)
		}

		// TODO: Actually validate the entry!

		if err := handler.CreateEntry(entry); err != nil {
			return StatusWithError(http.StatusInternalServerError, err)
		}
		// TODO: Have this return HTML or JSON or text/plain based on the header
		return Status(http.StatusOK)
	}))

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

type CustomHandler func(http.ResponseWriter, *http.Request) http.Handler

func (handler CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := handler(w, r); handler != nil {
		next.ServeHTTP(w, r)
		return
	}
}

func JSON(status int, v any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	})
}

// TODO: HTML functino similar to JSON, but have it use templ template?

// func Text(status int, content any) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//    w.Header().Set("Content-Type", "text/plain")
// 		w.WriteHeader(status)
// 		fmt.Fprint(w, content)
// 	})
// }

func StatusWithError(status int, err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, fmt.Sprintf("%s: %s", http.StatusText(status), err.Error()))
	})
}

func Status(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, http.StatusText(status))
	})
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if errors.Is(err, io.EOF) {
		return v, errors.New("body is empty or incomplete")
	}

	if err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
