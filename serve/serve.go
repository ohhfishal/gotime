package serve

import (
	"context"
	_ "embed"
	"log/slog"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
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
    {{ .error }}
  </h1>
</html>
`

func Serve(ctx context.Context, logger *slog.Logger, handler EntryHandler, port string) error {
	// TODO: Have this be configurable? Or just use env? Can probably handle this in cmd/serve.go
	// gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(loggingMiddleware(logger))
	r.Use(gin.Recovery())

	// TODO: Add this to assets and redirect
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/x-yaml", spec)
	})

	r.GET("/openapi/*filepath", gin.WrapH(http.StripPrefix("/openapi", swaggerui.Handler(spec))))
	r.GET("/assets/*filepath", gin.WrapH(http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets)))))
	r.GET("/favicon.ico", func(c *gin.Context) {
		http.ServeFileFS(c.Writer, c.Request, assets.Assets, "img/favicon.ico")
	})

	r.GET("/", func(c *gin.Context) {
		// TODO: Have this be handled via HTMX and GET endpoint
		entries, err := handler.GetAllEntries(entry.Today())
		if err != nil {
			// TODO: Make a template for this
			c.HTML(http.StatusInternalServerError, errorHTML, gin.H{"error": err.Error()})
			return
		}
		MainPage(
			entry.Summarize(entries),
		).Render(c, c.Writer)
	})

	api := r.Group("/api/v1")
	{
		api.GET("/entry", func(c *gin.Context) {
			entries, err := handler.GetAllEntries(entry.Today())
			if err != nil {
				c.HTML(http.StatusInternalServerError, errorHTML, gin.H{"error": err.Error()})
				return
			}
			Details(
				entry.Summarize(entries),
			).Render(c, c.Writer)

			// TODO: Filter entries based on some parameters
			// TODO: Have this return HTML or JSON based on the header
			// c.JSON(http.StatusOK, entries)
		})

		api.POST("/entry", func(c *gin.Context) {
			var entryData entry.Entry
			if err := c.ShouldBindJSON(&entryData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// TODO: Actually validate the entry!

			if err := handler.CreateEntry(entryData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// TODO: Have this return HTML or JSON or text/plain based on the header
			c.Status(http.StatusOK)
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
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

func loggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("replied to request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}
