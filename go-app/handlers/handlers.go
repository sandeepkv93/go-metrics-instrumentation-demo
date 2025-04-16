package handlers

import (
	"fmt"
	"net/http"

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
	// Increment request counter
	h.metrics.Requests.Add(r.Context(), 1)
	fmt.Fprintf(w, "Hello, World!")
}

// SetupRoutes configures all HTTP routes
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)
	return mux
}
