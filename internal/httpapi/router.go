package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"challenge-backend-arancia/internal/application/todos"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RouterOptions struct {
	TodoService *todos.Service
	Ready       func(ctx context.Context) error
	Logger      *slog.Logger
}

func NewRouter(opts RouterOptions) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestIDMiddleware())
	if opts.Logger != nil {
		r.Use(loggingMiddleware(opts.Logger))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		if opts.Ready == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()
		if err := opts.Ready(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	if opts.TodoService != nil {
		todoHandler{svc: opts.TodoService}.register(r)
	}

	return r
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const header = "X-Request-Id"
		reqID := c.GetHeader(header)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Writer.Header().Set(header, reqID)
		c.Set("request_id", reqID)
		c.Next()
	}
}

func loggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		reqID, _ := c.Get("request_id")

		logger.Info(
			"http_request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.Any("request_id", reqID),
		)
	}
}
