package tracing

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TracingProvider encapsulates the tracing setup
type TracingProvider struct {
	tracer     trace.Tracer
	provider   *sdktrace.TracerProvider
	propagator propagation.TextMapPropagator
}

// Config holds configuration for tracing setup
type Config struct {
	ServiceName       string
	ServiceVersion    string
	CollectorEndpoint string
	SampleRate        float64
}

// DefaultConfig provides sensible defaults
func DefaultConfig() Config {
	return Config{
		ServiceName:       "demo-service",
		ServiceVersion:    "1.0.0",
		CollectorEndpoint: "otel-collector:7317",
		SampleRate:        1.0, // Sample all traces for demonstration
	}
}

// NewTracingProvider initializes the tracing system
func NewTracingProvider(ctx context.Context, cfg Config) (*TracingProvider, error) {
	// Create gRPC connection to collector
	conn, err := grpc.DialContext(ctx, cfg.CollectorEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to create gRPC connection to collector: %v", err)
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Create OTLP exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("Failed to create trace exporter: %v", err)
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create and register tracer provider
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set up propagation
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Get a tracer
	tracer := tracerProvider.Tracer(cfg.ServiceName)

	tp := &TracingProvider{
		tracer:     tracer,
		provider:   tracerProvider,
		propagator: propagator,
	}

	log.Println("Tracing provider initialized successfully")
	return tp, nil
}

// Tracer returns the OpenTelemetry tracer
func (tp *TracingProvider) Tracer() trace.Tracer {
	return tp.tracer
}

// Propagator returns the context propagator
func (tp *TracingProvider) Propagator() propagation.TextMapPropagator {
	return tp.propagator
}

// Shutdown gracefully shuts down the tracing provider
func (tp *TracingProvider) Shutdown(ctx context.Context) error {
	log.Println("Shutting down tracer provider")
	return tp.provider.Shutdown(ctx)
}
