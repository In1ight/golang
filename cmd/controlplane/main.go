package main

import (
	"context"
	"errors"
	"example.com/analytics/internal/controlplane"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

func run() error {
	var ready atomic.Bool
	ready.Store(true)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", controlplane.HealthHandler)
	mux.HandleFunc("/readyz", controlplane.ReadyHandler(&ready))
	mux.Handle("/metrics", promhttp.Handler())

	cfg, err := controlplane.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var level slog.Level

	err = level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	logger.Info("Starting server", "addr", cfg.HTTPAddr)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("Shutdown signal received")
		ready.Store(false)

	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server stopped unexpectedly: %w", err)
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.ShutdownTimeout,
	)
	defer cancel()
	err = server.Shutdown(shutdownCtx)

	if err != nil {
		return fmt.Errorf("shutdown failed %w", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}
func main() {

	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
