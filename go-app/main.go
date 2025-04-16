package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	// Global meter
	meter = otel.GetMeterProvider().Meter("demo-service")
	// Counter instrument
	requests metric.Int64Counter
)

func initMeter() func() {
	// Create OTLP exporter with custom port
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint("otel-collector:7317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	// Create resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("demo-service"),
	)

	// Create meter provider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(1*time.Second)),
		),
	)

	// Set global meter provider
	otel.SetMeterProvider(provider)

	// Create our metrics instruments
	var err2 error
	requests, err2 = meter.Int64Counter(
		"demo_requests_total",
		metric.WithDescription("Total number of requests received"),
	)
	if err2 != nil {
		log.Fatalf("Failed to create counter: %v", err2)
	}

	// Return a function to shutdown the exporter when the application exits
	return func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Increment the counter
	requests.Add(context.Background(), 1)
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Initialize OpenTelemetry
	shutdown := initMeter()
	defer shutdown()

	http.HandleFunc("/", handler)

	fmt.Println("Serving on :8085...")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal(err)
	}
}
