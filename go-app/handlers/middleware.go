package handlers

import (
	"net/http"

	"github.com/example/go-app/tracing"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware wraps an HTTP handler with request tracing
func TracingMiddleware(next http.Handler, tp *tracing.TracingProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract context from HTTP headers
		ctx := r.Context()
		propagator := tp.Propagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		// Start a span
		operation := r.URL.Path
		if operation == "" {
			operation = "/"
		}

		spanName := "HTTP " + r.Method + " " + operation
		ctx, span := tp.Tracer().Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPTargetKey.String(r.URL.Path),
				semconv.HTTPURLKey.String(r.URL.String()),
				semconv.HTTPUserAgentKey.String(r.UserAgent()),
				semconv.HTTPSchemeKey.String(r.URL.Scheme),
				semconv.NetHostNameKey.String(r.Host),
			),
		)
		defer span.End()

		// Inject the context back into the request
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
