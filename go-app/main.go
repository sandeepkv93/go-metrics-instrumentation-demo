package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-app/handlers"
	"github.com/example/go-app/metrics"
	"github.com/example/go-app/tracing"
)

func main() {
	// Create a context for initialization
	ctx := context.Background()

	// Initialize metrics with default config
	metricsProvider, err := metrics.NewMetricsProvider(ctx, metrics.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer metricsProvider.Shutdown(context.Background())

	// Initialize tracing with default config
	tracingProvider, err := tracing.NewTracingProvider(ctx, tracing.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer tracingProvider.Shutdown(context.Background())

	// Initialize HTTP handlers with both providers
	handler := handlers.NewHandler(metricsProvider, tracingProvider)

	// Create server
	server := &http.Server{
		Addr:    ":8085",
		Handler: handler.SetupRoutes(),
	}

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server running on :8085...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-done

	// Graceful shutdown
	log.Println("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
