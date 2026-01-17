package main

import (
	"context"
	"errors"
	"finance-manager/internal/data"
	"finance-manager/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
	srv := service.NewService(db, logger)
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
	go func() {
		logger.Info("listening on", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen failed", "err", err)
			os.Exit(1) // Exit if the server fails to start
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "err", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
