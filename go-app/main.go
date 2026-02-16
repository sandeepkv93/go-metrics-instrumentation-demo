package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-app/handlers"
	"github.com/example/go-app/logging"
	"github.com/example/go-app/metrics"
	"github.com/example/go-app/tracing"
)

func main() {
	ctx := context.Background()

	loggingProvider, err := logging.NewLogProvider(ctx, logging.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}
	defer loggingProvider.Shutdown(context.Background())
	logger := loggingProvider.Logger()

	metricsProvider, err := metrics.NewMetricsProvider(ctx, metrics.DefaultConfig())
	if err != nil {
		logger.Error("Failed to initialize metrics", "error", err)
		os.Exit(1)
	}
	defer metricsProvider.Shutdown(context.Background())

	tracingProvider, err := tracing.NewTracingProvider(ctx, tracing.DefaultConfig())
	if err != nil {
		logger.Error("Failed to initialize tracing", "error", err)
		os.Exit(1)
	}
	defer tracingProvider.Shutdown(context.Background())

	handler := handlers.NewHandler(metricsProvider, tracingProvider, loggingProvider)

	server := &http.Server{
		Addr:    ":8085",
		Handler: handler.SetupRoutes(),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server running", "address", ":8085")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	<-done

	logger.Info("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}
