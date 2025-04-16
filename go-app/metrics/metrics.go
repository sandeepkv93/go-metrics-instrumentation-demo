package metrics

import (
	"context"
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
	Requests metric.Int64Counter
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
		return nil, err
	}

	// Create resource with service information
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
	)

	// Create meter provider
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(cfg.ExportInterval)),
		),
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

	// Initialize metrics
	requests, err := meter.Int64Counter(
		"demo_requests_total",
		metric.WithDescription("Total number of requests received"),
	)
	if err != nil {
		return nil, err
	}
	mp.Requests = requests

	return mp, nil
}

// Shutdown gracefully shuts down the metrics provider
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	return mp.provider.Shutdown(ctx)
}
