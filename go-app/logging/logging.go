package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

// LogProvider encapsulates logging setup
type LogProvider struct {
	logger zerolog.Logger
}

// Config holds configuration for logging setup
type Config struct {
	ServiceName      string
	LogLevel         string
	LogFilePath      string
	ConsoleLoggingOn bool
	JSONFormatting   bool
}

// DefaultConfig provides sensible defaults
func DefaultConfig() Config {
	return Config{
		ServiceName:      "demo-service",
		LogLevel:         "info",
		LogFilePath:      "/var/log/go-app/app.log",
		ConsoleLoggingOn: true,
		JSONFormatting:   true,
	}
}

// NewLogProvider initializes the logging system
func NewLogProvider(cfg Config) (*LogProvider, error) {
	// Create log directory if it doesn't exist
	logDir := "/var/log/go-app"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create or open log file
	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Set global level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure timestamp format
	zerolog.TimeFieldFormat = time.RFC3339

	// Configure output (file, console, or both)
	var writers []io.Writer
	writers = append(writers, logFile)

	if cfg.ConsoleLoggingOn {
		if cfg.JSONFormatting {
			writers = append(writers, os.Stdout)
		} else {
			// Use console writer for pretty output during development
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
			writers = append(writers, consoleWriter)
		}
	}

	// Create multi-writer
	multiWriter := zerolog.MultiLevelWriter(writers...)

	// Create logger with service context
	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Logger()

	// Set as default logger
	log.Logger = logger

	fmt.Println("Logging provider initialized successfully")
	return &LogProvider{
		logger: logger,
	}, nil
}

// Logger returns the zerolog logger
func (lp *LogProvider) Logger() zerolog.Logger {
	return lp.logger
}

// ContextLogger returns a logger with context values
func (lp *LogProvider) ContextLogger(ctx context.Context) zerolog.Logger {
	traceID := TraceIDFromContext(ctx)
	spanID := SpanIDFromContext(ctx)

	logger := lp.logger

	// Add trace context if available
	if traceID != "" {
		// IMPORTANT: Use "trace_id" exactly as this is what Grafana expects
		logger = logger.With().
			Str("trace_id", traceID).
			Str("span_id", spanID).
			Logger()
	}

	return logger
}

// ShutDown properly flushes and closes log resources
func (lp *LogProvider) Shutdown(ctx context.Context) error {
	// No specific shutdown needed for zerolog
	return nil
}

// SpanIDFromContext extracts span ID from OpenTelemetry context
func SpanIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.SpanID().String()
}

// TraceIDFromContext extracts trace ID from OpenTelemetry context
func TraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	// Format exactly as Tempo expects - lowercase hex without dashes
	return spanCtx.TraceID().String()
}
