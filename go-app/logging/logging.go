package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	otlploggrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// LogProvider encapsulates logging setup.
type LogProvider struct {
	logger   *slog.Logger
	provider *sdklog.LoggerProvider
}

// Config holds configuration for logging setup.
type Config struct {
	ServiceName       string
	LogLevel          string
	CollectorEndpoint string
	OTLPInsecure      bool
	OTELLogsEnabled   bool
}

// DefaultConfig provides sensible defaults.
func DefaultConfig() Config {
	return Config{
		ServiceName:       "demo-service",
		LogLevel:          "info",
		CollectorEndpoint: "otel-collector:7317",
		OTLPInsecure:      true,
		OTELLogsEnabled:   true,
	}
}

type multiHandler struct {
	handlers []slog.Handler
}

type traceContextHandler struct {
	next slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		next = append(next, handler.WithAttrs(attrs))
	}
	return &multiHandler{handlers: next}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		next = append(next, handler.WithGroup(name))
	}
	return &multiHandler{handlers: next}
}

func (h *traceContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *traceContextHandler) Handle(ctx context.Context, r slog.Record) error {
	traceID := ""
	spanID := ""
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		traceID = sc.TraceID().String()
		spanID = sc.SpanID().String()
	}
	r.AddAttrs(
		slog.String("trace_id", traceID),
		slog.String("span_id", spanID),
	)
	return h.next.Handle(ctx, r)
}

func (h *traceContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceContextHandler{next: h.next.WithAttrs(attrs)}
}

func (h *traceContextHandler) WithGroup(name string) slog.Handler {
	return &traceContextHandler{next: h.next.WithGroup(name)}
}

var (
	loggerMu     sync.RWMutex
	globalLogger *slog.Logger
)

// NewLogProvider initializes logging with stdout JSON and optional OTLP export.
func NewLogProvider(ctx context.Context, cfg Config) (*LogProvider, error) {
	level := parseLogLevel(cfg.LogLevel)
	stdout := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	if !cfg.OTELLogsEnabled {
		logger := slog.New(&traceContextHandler{next: stdout}).With("service", cfg.ServiceName)
		loggerMu.Lock()
		globalLogger = logger
		loggerMu.Unlock()
		slog.SetDefault(logger)
		logger.Info("logging initialized", "otlp_logs_enabled", false)
		return &LogProvider{logger: logger}, nil
	}

	opts := []otlploggrpc.Option{otlploggrpc.WithEndpoint(cfg.CollectorEndpoint)}
	if cfg.OTLPInsecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}
	exporter, err := otlploggrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create otlp log exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create logs resource: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	otelHandler := otelslog.NewHandler(cfg.ServiceName, otelslog.WithLoggerProvider(lp))
	logger := slog.New(&traceContextHandler{next: &multiHandler{handlers: []slog.Handler{stdout, otelHandler}}}).With("service", cfg.ServiceName)

	loggerMu.Lock()
	globalLogger = logger
	loggerMu.Unlock()
	slog.SetDefault(logger)

	logger.Info("logging initialized", "otlp_logs_enabled", true, "collector_endpoint", cfg.CollectorEndpoint)
	return &LogProvider{logger: logger, provider: lp}, nil
}

// Logger returns the slog logger.
func (lp *LogProvider) Logger() *slog.Logger {
	if lp != nil && lp.logger != nil {
		return lp.logger
	}
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	if globalLogger != nil {
		return globalLogger
	}
	return slog.Default()
}

// ContextLogger returns logger instance. Trace/span enrichment happens at emit time via context.
func (lp *LogProvider) ContextLogger(ctx context.Context) *slog.Logger {
	return lp.Logger()
}

// Shutdown flushes OTLP log resources.
func (lp *LogProvider) Shutdown(ctx context.Context) error {
	if lp == nil || lp.provider == nil {
		return nil
	}
	return lp.provider.Shutdown(ctx)
}

func parseLogLevel(v string) slog.Level {
	switch v {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
