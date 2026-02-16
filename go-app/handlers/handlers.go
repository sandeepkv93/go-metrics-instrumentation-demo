package handlers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/example/go-app/logging"
	"github.com/example/go-app/metrics"
	"github.com/example/go-app/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Handler encapsulates HTTP handlers and their dependencies
// for metrics, tracing and logging.
type Handler struct {
	metrics *metrics.MetricsProvider
	tracing *tracing.TracingProvider
	logging *logging.LogProvider
}

// NewHandler creates a new handler with dependencies.
func NewHandler(metrics *metrics.MetricsProvider, tracing *tracing.TracingProvider, logging *logging.LogProvider) *Handler {
	return &Handler{
		metrics: metrics,
		tracing: tracing,
		logging: logging,
	}
}

// Home handles the root path.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logging.ContextLogger(ctx)

	ctx, span := h.tracing.Tracer().Start(ctx, "home_handler_logic")
	defer span.End()

	logger.InfoContext(ctx,
		"Processing home request",
		"http.method", r.Method,
		"http.url", r.URL.String(),
		"http.user_agent", r.UserAgent(),
	)

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.user_agent", r.UserAgent()),
	)

	span.AddEvent("request_received", trace.WithAttributes(
		attribute.String("remote_addr", r.RemoteAddr),
	))

	start := time.Now()
	h.metrics.Requests.Add(ctx, 1)

	if err := h.validateRequest(ctx, r); err != nil {
		logger.ErrorContext(ctx, "Request validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		span.SetStatus(codes.Error, err.Error())
		span.AddEvent("request_validation_failed")
		return
	}

	userData, err := h.fetchUserData(ctx)
	if err != nil {
		span.AddEvent("user_data_fetch_failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
	}

	weatherData, err := h.fetchWeatherData(ctx)
	if err != nil {
		span.AddEvent("weather_data_fetch_failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
	}

	result, err := h.processBusinessLogic(ctx, userData, weatherData)
	if err != nil {
		http.Error(w, "Processing error", http.StatusInternalServerError)
		span.SetStatus(codes.Error, "Business logic processing failed")
		span.RecordError(err)
		return
	}

	if err := h.updateCache(ctx, result); err != nil {
		span.RecordError(err)
		span.AddEvent("cache_update_failed")
	}

	responseData, err := h.renderResponse(ctx, result)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		span.SetStatus(codes.Error, "Response rendering failed")
		span.RecordError(err)
		return
	}

	fmt.Fprint(w, responseData)
	h.metrics.ObserveRequestDuration(ctx, start)

	duration := time.Since(start)
	span.SetAttributes(attribute.Int64("response.time_ms", duration.Milliseconds()))
	span.AddEvent("request_completed")

	logger.InfoContext(ctx,
		"Request completed successfully",
		"duration_ms", float64(duration.Microseconds())/1000.0,
	)
}

// validateRequest simulates request validation.
func (h *Handler) validateRequest(ctx context.Context, r *http.Request) error {
	ctx, span := h.tracing.Tracer().Start(ctx, "validate_request")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	span.SetAttributes(attribute.String("http.method", r.Method))
	time.Sleep(time.Duration(10+rand.Intn(30)) * time.Millisecond)

	if rand.Intn(100) < 5 {
		err := errors.New("invalid request parameters")
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "")

	logger.InfoContext(ctx,
		"Request validation completed successfully",
		"http.method", r.Method,
		"http.url", r.URL.String(),
		"http.user_agent", r.UserAgent(),
	)
	return nil
}

// fetchUserData simulates a database call.
func (h *Handler) fetchUserData(ctx context.Context) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "fetch_user_data")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	span.SetAttributes(
		attribute.String("db.system", "postgres"),
		attribute.String("db.operation", "SELECT"),
	)

	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	if rand.Intn(100) < 10 {
		err := errors.New("database connection error")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", err
	}

	userData := fmt.Sprintf("user-%d", rand.Intn(1000))
	span.SetAttributes(attribute.String("user.id", userData))
	span.SetStatus(codes.Ok, "")

	logger.InfoContext(ctx,
		"User data fetched successfully",
		"db.system", "postgres",
		"db.operation", "SELECT",
	)
	return userData, nil
}

// fetchWeatherData simulates an external API call.
func (h *Handler) fetchWeatherData(ctx context.Context) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "fetch_weather_data")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	span.SetAttributes(
		attribute.String("http.url", "https://api.weather.example.com/current"),
		attribute.String("http.method", "GET"),
	)

	time.Sleep(time.Duration(100+rand.Intn(150)) * time.Millisecond)

	if rand.Intn(100) < 15 {
		err := errors.New("weather API timeout")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", err
	}

	temps := []string{"sunny", "cloudy", "rainy", "snowy"}
	weather := temps[rand.Intn(len(temps))]
	span.SetAttributes(attribute.String("weather.condition", weather))
	span.SetStatus(codes.Ok, "")

	logger.InfoContext(ctx,
		"Weather data fetched successfully",
		"http.method", "GET",
		"http.url", "https://api.weather.example.com/current",
	)
	return weather, nil
}

// processBusinessLogic simulates complex business logic.
func (h *Handler) processBusinessLogic(ctx context.Context, userData, weatherData string) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "process_business_logic")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	span.SetAttributes(
		attribute.String("user.data", userData),
		attribute.String("weather.data", weatherData),
	)

	result, err := h.processStep1(ctx, userData)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	result, err = h.processStep2(ctx, result, weatherData)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetStatus(codes.Ok, "")
	logger.InfoContext(ctx,
		"Business logic processed successfully",
		"process_business_logic_trace_id", span.SpanContext().TraceID().String(),
	)
	return result, nil
}

// processStep1 is a sub-step of business logic.
func (h *Handler) processStep1(ctx context.Context, userData string) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "process_step1")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	time.Sleep(time.Duration(20+rand.Intn(40)) * time.Millisecond)

	result := fmt.Sprintf("processed-%s", userData)
	span.SetAttributes(attribute.String("step1.result", result))

	logger.InfoContext(ctx,
		"Step 1 processed successfully",
		"process_step1_trace_id", span.SpanContext().TraceID().String(),
	)
	return result, nil
}

// processStep2 is a sub-step of business logic.
func (h *Handler) processStep2(ctx context.Context, step1Result, weatherData string) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "process_step2")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	time.Sleep(time.Duration(30+rand.Intn(50)) * time.Millisecond)

	if rand.Intn(100) < 5 {
		err := errors.New("processing algorithm error")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", err
	}

	result := fmt.Sprintf("%s-with-%s", step1Result, weatherData)
	span.SetAttributes(attribute.String("step2.result", result))

	logger.InfoContext(ctx,
		"Step 2 processed successfully",
		"process_step2_trace_id", span.SpanContext().TraceID().String(),
	)
	return result, nil
}

// updateCache simulates writing to a cache.
func (h *Handler) updateCache(ctx context.Context, data string) error {
	ctx, span := h.tracing.Tracer().Start(ctx, "update_cache")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	span.SetAttributes(
		attribute.String("cache.system", "redis"),
		attribute.String("cache.operation", "SET"),
		attribute.String("cache.key", "result"),
	)

	time.Sleep(time.Duration(5+rand.Intn(20)) * time.Millisecond)

	if rand.Intn(100) < 8 {
		err := errors.New("cache connection error")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "")
	logger.InfoContext(ctx,
		"Cache updated successfully",
		"cache.system", "redis",
		"cache.operation", "SET",
		"cache.key", "result",
	)
	return nil
}

// renderResponse simulates rendering the response.
func (h *Handler) renderResponse(ctx context.Context, data string) (string, error) {
	ctx, span := h.tracing.Tracer().Start(ctx, "render_response")
	defer span.End()

	logger := h.logging.ContextLogger(ctx)

	time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)

	response := fmt.Sprintf("Hello, World! Result: %s", data)
	span.SetAttributes(attribute.Int("response.size", len(response)))

	logger.InfoContext(ctx,
		"Response rendered successfully",
		"render_response_trace_id", span.SpanContext().TraceID().String(),
	)
	return response, nil
}

// SetupRoutes configures all HTTP routes.
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)

	// Wrap everything with tracing middleware.
	return TracingMiddleware(mux, h.tracing)
}
