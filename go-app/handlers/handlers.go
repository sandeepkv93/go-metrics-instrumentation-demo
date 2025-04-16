package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/example/go-app/metrics"
)

// Handler encapsulates HTTP handlers and their dependencies
type Handler struct {
	metrics *metrics.MetricsProvider
}

// NewHandler creates a new handler with dependencies
func NewHandler(metrics *metrics.MetricsProvider) *Handler {
	return &Handler{
		metrics: metrics,
	}
}

// Home handles the root path
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Start measuring request duration
	start := time.Now()

	// Increment request counter
	h.metrics.Requests.Add(r.Context(), 1)

	// Simulate variable work (between 50-250ms)
	sleepTime := 50 + rand.Intn(200)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	fmt.Fprintf(w, "Hello, World! (Response took %dms)", sleepTime)

	// Record request duration
	h.metrics.ObserveRequestDuration(r.Context(), start)
}

// SetupRoutes configures all HTTP routes
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)
	return mux
}
