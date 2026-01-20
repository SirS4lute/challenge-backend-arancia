package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"challenge-backend-arancia/internal/application/todos"
	"challenge-backend-arancia/internal/config"
	"challenge-backend-arancia/internal/httpapi"
	"challenge-backend-arancia/internal/storage/boltdb"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.FromEnv()

	gin.SetMode(cfg.GinMode)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	db, err := boltdb.Open(cfg.DBPath)
	if err != nil {
		panic(err)
	}
	defer func() { _ = db.Close() }()

	repo, err := boltdb.NewTodoRepository(db)
	if err != nil {
		panic(err)
	}
	svc, err := todos.NewService(repo, todos.UUIDGenerator{})
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Port),
		Handler: httpapi.NewRouter(httpapi.RouterOptions{
			TodoService: svc,
			Ready: func(ctx context.Context) error {
				_, err := repo.List(ctx)
				return err
			},
			Logger: logger,
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func parseLogLevel(v string) slog.Level {
	switch v {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		log.Printf("unknown LOG_LEVEL=%q, defaulting to info", v)
		return slog.LevelInfo
	}
}

