# Go Observability Instrumentation and Visualization Demo 📊 📈 🔍

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## 🔍 Overview

This repository demonstrates a modern approach to instrumenting Go applications with both metrics and traces using OpenTelemetry, showing how to implement a complete observability stack with minimal code changes.

## 🌟 Features

- ✅ Complete working examples of metrics and tracing instrumentation
- ✅ Docker Compose setup for quick deployment
- ✅ Pre-configured Grafana dashboards
- ✅ Request count and latency histogram metrics
- ✅ Distributed tracing with detailed span information
- ✅ Correlation between metrics and traces in a single UI

## 🔄 Branches

### 🚀 otel-tempo-mimir-grafana (default)

**Modern OpenTelemetry Approach with Complete Observability**

This branch demonstrates the modern approach to observability using the OpenTelemetry standard:

- **Instrumentation**: Uses OpenTelemetry SDK for Go (metrics and traces)
- **Collection**: OpenTelemetry Collector
- **Metrics Storage**: Grafana Mimir (Prometheus-compatible)
- **Trace Storage**: Grafana Tempo
- **Visualization**: Grafana dashboards and Explore

**Data Flow**:

- **Metrics**: Go App → OTel SDK → OTel Collector → Mimir → Grafana
- **Traces**: Go App → OTel SDK → OTel Collector → Tempo → Grafana

### 🏛️ prometheus-and-mimir

**Traditional Prometheus Approach**

This branch demonstrates the classic Prometheus instrumentation approach:

- **Instrumentation**: Prometheus client library for Go
- **Collection**: Direct scraping by Prometheus
- **Storage**: Grafana Mimir
- **Visualization**: Same Grafana dashboards (for direct comparison)

**Data Flow**: Go App → Prometheus client → Prometheus → Mimir → Grafana

## 🛠️ Setup

### Prerequisites

- 🐳 Docker and Docker Compose
- 🧪 Git
- 🔄 curl (for testing)

### Quick Start

1. Clone the repository:

   ```bash
   git clone git@github.com:sandeepkv93/go-metrics-instrumentation-demo.git
   cd go-metrics-instrumentation-demo
   ```

2. Start the stack:

   ```bash
   docker-compose up -d
   ```

3. Access the Grafana dashboard:

   - URL: http://localhost:3000
   - Default user/pass: admin/admin

4. Generate some traffic:

   ```bash
   ./hit-loop.sh
   ```

   # Make it executable

   chmod +x hit-loop.sh

   # Run it

   ```
   ./hit-loop.sh 200 0.2
   ```

## 📊 Exploring Metrics and Traces

### Metrics

Navigate to the pre-configured dashboard in Grafana to see:

- Request counts
- Request duration percentiles
- Error rates

### Traces

In Grafana, use the Explore tab and select the Tempo data source to:

1. View all traces
2. Search by trace ID
3. Filter by service, duration, or status

### Correlation

From a metrics dashboard, you can:

1. Select a time range with interesting patterns
2. Click "Explore" to dive into traces from that time period
3. Find specific traces that explain metric anomalies

## 📚 Code Structure

- `/go-app/metrics` - Metrics instrumentation code
- `/go-app/tracing` - Tracing instrumentation code
- `/go-app/handlers` - HTTP handlers with observability
- `/grafana` - Grafana configuration and dashboards
- `/tempo-config.yaml` - Tempo configuration
- `/otel-collector-config.yaml` - OpenTelemetry Collector configuration

## 📊 Included Metrics

- **Request Count**: Total number of requests
- **Request Duration**: Histogram of request durations
- **Request Duration Percentiles**: p50, p90, p95, and p99 latency
- **Request Rate**: Requests per second
- **Average Latency**: Average request duration

## 🧩 Architecture

### 🚀 OpenTelemetry Stack (default branch)

```
┌────────────┐    ┌─────────────────┐    ┌──────────────────┐    ┌─────────┐
│            │    │                 │    │                  │    │         │
│  Go App    ├───►│  OTel Collector ├─── │  Grafana Mimir   ├───►│ Grafana │
│ (OTel SDK) │    │                 │    │                  │    │         │
└────────────┘    └─────────────────┘    └──────────────────┘    └─────────┘
```

### 🏛️ Prometheus Stack (prometheus-and-mimir branch)

```
┌─────────────┐    ┌─────────────┐    ┌──────────────────┐    ┌─────────┐
│             │    │             │    │                  │    │         │
│  Go App     │◄───┤  Prometheus ├───►│  Grafana Mimir   ├───►│ Grafana │
│(Prom client)│    │             │    │                  │    │         │
└─────────────┘    └─────────────┘    └──────────────────┘    └─────────┘
```

## 📝 Configuration Files

- **docker-compose.yml**: Docker Compose configuration
- **otel-collector-config.yaml**: OpenTelemetry Collector configuration
- **mimir-config.yaml**: Mimir configuration
- **grafana/provisioning/**: Grafana provisioning files
- **hit-loop.sh**: Script to generate traffic

## 🔍 Exploring Metrics

Use Grafana's Explore feature to explore all available metrics:

1. Open Grafana (http://localhost:3000)
2. Click on the Explore icon in the left sidebar
3. Try queries like:
   - `{__name__=~".*request.*"}`
   - `rate(demo_requests_total[1m])`
   - `histogram_quantile(0.95, sum(rate(demo_request_duration_seconds_bucket[1m])) by (le))`

## 📊 Example PromQL Queries

Here are some useful queries for analyzing your application:

```
# Request Rate (per second)
rate(demo_requests_total[1m])

# 95th Percentile Latency
histogram_quantile(0.95, sum(rate(demo_request_duration_seconds_bucket[1m])) by (le))

# Error Rate
rate(demo_request_errors_total[1m])

# Request Rate by Endpoint
sum(rate(demo_requests_total[1m])) by (endpoint)
```

## 🤔 Why Two Approaches?

- **OpenTelemetry**: Newer standard that provides a unified approach for metrics, traces, and logs
- **Prometheus**: Well-established standard for metrics with a large ecosystem
- **Comparison**: Helps you decide which approach fits your needs better

## 🔄 Switching Branches

To switch between the two implementations:

```bash
# Stop the current stack
docker-compose down

# Switch branch
git checkout prometheus-and-mimir  # or otel-mimir-grafana

# Start the stack again
docker-compose up -d
```

## 🔧 Troubleshooting

### Common Issues

1. **OTel Collector Not Starting**:

   ```bash
   # Check logs
   docker-compose logs otel-collector

   # Validate configuration
   docker run --rm -v $(pwd)/otel-collector-config.yaml:/config.yaml otel/opentelemetry-collector validate --config=/config.yaml
   ```

2. **Metrics Not Showing in Grafana**:

   ```bash
   # Check Mimir targets
   curl http://localhost:9009/api/v1/targets | jq

   # Verify metrics are being collected
   curl http://localhost:9009/api/v1/query?query=up
   ```

3. **Docker Compose Network Issues**:
   ```bash
   # Recreate network
   docker-compose down
   docker network prune -f
   docker-compose up -d
   ```

- **Go Application**: Simple web service instrumented with OpenTelemetry
- **OpenTelemetry Collector**: Receives, processes, and exports telemetry data
- **Mimir**: Stores and queries time-series metrics
- **Tempo**: Stores and queries distributed traces
- **Grafana**: Visualizes metrics and traces in a unified UI

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- OpenTelemetry community
- Prometheus community
- Grafana Labs for Mimir and Grafana

---

Made with ❤️ by [Sandeep Vishnu](https://github.com/sandeepkv93)
