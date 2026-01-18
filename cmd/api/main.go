package main

import (
	"context"
	"errors"
	"finance-manager/internal"
	"finance-manager/internal/data"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not loaded, using system environment variables")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("DATABASE_URL environment variable required")
		os.Exit(1)
	}
	dbPool, err := data.NewDB(dsn)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	logger.Info("Successfully connected to the Finance Manager database!")

	db := data.NewModels(dbPool)
	srv := internal.NewService(db, logger)
	api := NewAPI(srv, logger)
	handler := api.Router()

	server := &http.Server{
		Addr:         ":4000",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("listening on", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-quit:
		logger.Info("shutting down server...")
	case err := <-serverErr:
		logger.Error("server listen failed", "err", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "err", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}
