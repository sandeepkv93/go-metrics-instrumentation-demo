package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// MetricsProvider encapsulates the metrics setup
type MetricsProvider struct {
	meter    metric.Meter
	provider *sdkmetric.MeterProvider
	// Expose instrumented metrics
	Requests        metric.Int64Counter
	RequestDuration metric.Float64Histogram
}

// Config holds configuration for metrics setup
type Config struct {
	ServiceName       string
	CollectorEndpoint string
	ExportInterval    time.Duration
}

// DefaultConfig provides sensible defaults
func DefaultConfig() Config {
	return Config{
		ServiceName:       "demo-service",
		CollectorEndpoint: "otel-collector:7317",
		ExportInterval:    1 * time.Second,
	}
}

// NewMetricsProvider initializes the metrics system
func NewMetricsProvider(ctx context.Context, cfg Config) (*MetricsProvider, error) {
	// Create OTLP exporter
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("Failed to create exporter: %v", err)
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create resource with service information
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	)

	// Define histogram boundaries explicitly
	histogramBoundaries := []float64{
		0.001, 0.002, 0.005, 0.01, 0.02, 0.05,
		0.1, 0.2, 0.3, 0.4, 0.5, 0.75, 1, 1.5, 2, 5,
	}

	// Create view to customize histogram buckets
	view := sdkmetric.NewView(
		sdkmetric.Instrument{
			Name: "demo_request_duration_seconds",
			Kind: sdkmetric.InstrumentKindHistogram,
		},
		sdkmetric.Stream{
			Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: histogramBoundaries,
			},
		},
	)

	// Create reader with explicit debug logging
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(cfg.ExportInterval),
	)

	// Create meter provider with the view
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(view),
	)

	// Set global meter provider
	otel.SetMeterProvider(provider)

	// Create meter
	meter := provider.Meter(cfg.ServiceName)

	// Create metrics provider
	mp := &MetricsProvider{
		meter:    meter,
		provider: provider,
	}

	// Initialize request counter metric
	requests, err := meter.Int64Counter(
		"demo_requests_total",
		metric.WithDescription("Total number of requests received"),
	)
	if err != nil {
		log.Printf("Failed to create counter: %v", err)
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}
	mp.Requests = requests

	// Initialize request duration histogram metric
	requestDuration, err := meter.Float64Histogram(
		"demo_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Printf("Failed to create histogram: %v", err)
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}
	mp.RequestDuration = requestDuration

	log.Println("Metrics provider initialized successfully")
	return mp, nil
}

// ObserveRequestDuration is a helper to measure the duration of a request
func (mp *MetricsProvider) ObserveRequestDuration(ctx context.Context, start time.Time) {
	duration := time.Since(start).Seconds()
	log.Printf("Recording request duration: %.6f seconds", duration)
	mp.RequestDuration.Record(ctx, duration)
}

// Shutdown gracefully shuts down the metrics provider
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	return mp.provider.Shutdown(ctx)
}
