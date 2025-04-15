package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "demo_requests_total",
			Help: "Total number of requests received",
		},
	)
)

func handler(w http.ResponseWriter, r *http.Request) {
	requests.Inc()
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	prometheus.MustRegister(requests)

	http.Handle("/", http.HandlerFunc(handler))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Serving on :8085...")
	http.ListenAndServe(":8085", nil)
}
