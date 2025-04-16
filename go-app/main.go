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
	// Create a context for initialization
	ctx := context.Background()

	// Initialize logging with default config
	loggingProvider, err := logging.NewLogProvider(logging.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}
	logger := loggingProvider.Logger()

	// Initialize metrics with default config
	metricsProvider, err := metrics.NewMetricsProvider(ctx, metrics.DefaultConfig())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize metrics")
	}
	defer metricsProvider.Shutdown(context.Background())

	// Initialize tracing with default config
	tracingProvider, err := tracing.NewTracingProvider(ctx, tracing.DefaultConfig())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize tracing")
	}
	defer tracingProvider.Shutdown(context.Background())

	// Initialize HTTP handlers with all providers
	handler := handlers.NewHandler(metricsProvider, tracingProvider, loggingProvider)

	// Create server
	server := &http.Server{
		Addr:    ":8085",
		Handler: handler.SetupRoutes(),
	}

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info().Str("address", ":8085").Msg("Server running")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server error")
		}
	}()

	// Wait for termination signal
	<-done

	// Graceful shutdown
	logger.Info().Msg("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited gracefully")
}
