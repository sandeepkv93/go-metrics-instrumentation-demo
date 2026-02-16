# Go Observability Instrumentation and Visualization Demo

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## Overview

This repository demonstrates a modern approach to instrumenting Go applications with logs, metrics and traces using OpenTelemetry, showing how to implement a complete observability stack with minimal code changes.

## Features

- Complete working examples of logs, metrics and tracing instrumentation
- Docker Compose setup for quick deployment
- Pre-configured Grafana dashboards
- Request count and latency histogram metrics
- Distributed tracing with detailed span information
- Correlation between metrics, traces and logs in a single UI

## Branches

### main (default)

**Current primary branch: complete tri-signal OpenTelemetry stack (metrics, traces, logs)**

- **Metrics**: Go App → OTel SDK → OTel Collector → Mimir → Grafana
- **Traces**: Go App → OTel SDK → OTel Collector → Tempo → Grafana
- **Logs**: Go App → slog/OTel Log SDK → OTel Collector → Loki → Grafana

### Other branches

- `otel-loki-tempo-mimir-grafana`: earlier full observability branch (historically used Zerolog + Promtail for logs)
- `otel-loki-tempo-mimir-grafana-no-promtail`: variant branch for log-pipeline experimentation
- `otel-promtail-loki-tempo-mimir-grafana`: explicit Promtail-based logs branch
- `otel-tempo-mimir-grafana`: OpenTelemetry metrics + traces only (no logs pipeline)
- `otel-mimir-grafana`: metrics-focused OpenTelemetry branch
- `prometheus-and-mimir`: traditional Prometheus instrumentation branch

## Setup

### Prerequisites

- Docker and Docker Compose
- Git
- curl (for testing)

### Quick Start

1. Clone the repository:

   ```bash
   git clone git@github.com:sandeepkv93/go-observability-demo.git
   cd go-observability-demo
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

## Exploring Metrics and Traces

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

### Logs

Use Grafana's Explore tab with the Loki data source to:

1. View all logs with `{job="go-app"}`
2. Filter by service with `{service="demo-service"}`
3. Filter by log level with `{job="go-app", level="error"}`
4. Search for specific text with `{job="go-app"} |= "error"`

### Correlation

This demo features full correlation between metrics, traces, and logs:

1. From a metric dashboard, select a time range and click "Explore" to find traces
2. From a trace view, click "Logs for this span" to see relevant logs
3. From logs, click on a trace ID to jump to the corresponding trace

All components automatically share context:

- Trace IDs are included in logs via OpenTelemetry context propagation
- Service names are consistently used across all telemetry types
- Timestamps are synchronized for temporal correlation

## Code Structure

### Application Code

- `/go-app/metrics` - Metrics instrumentation code
- `/go-app/tracing` - Tracing instrumentation code
- `/go-app/logging` - Logging instrumentation with slog + OpenTelemetry logs bridge
- `/go-app/handlers` - HTTP handlers with observability

### Configuration

- `/configs` - All configuration files for observability tools
  - `/configs/grafana` - Grafana configuration and dashboards
    - `/configs/grafana/provisioning` - Datasources and dashboards
  - `/configs/loki` - Loki log aggregation configuration
  - `/configs/mimir` - Mimir metrics configuration
  - `/configs/tempo` - Tempo tracing configuration
  - `/configs/otel-collector` - OpenTelemetry Collector configuration

### Infrastructure

- `docker-compose.yml` - Container orchestration
- `hit-loop.sh` - Test script for generating traffic

## Included Telemetry

### Metrics

- **Request Count**: Total number of requests
- **Request Duration**: Histogram of request durations
- **Request Duration Percentiles**: p50, p90, p95, and p99 latency
- **Request Rate**: Requests per second
- **Average Latency**: Average request duration

### Logs

- **Structured JSON Logs**: All logs are in structured JSON format
- **Context-Enriched**: Logs include trace IDs, span IDs, and service name
- **Level-Based Filtering**: Support for filtering by log level (info, error, etc.)

## Architecture

### main (default)

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Tempo -> Grafana
```

#### Logs Flow
```text
Go App (slog + OTel Log SDK) -> OTel Collector -> Grafana Loki -> Grafana
```

### otel-loki-tempo-mimir-grafana

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Tempo -> Grafana
```

#### Logs Flow
```text
Go App (Zerolog) -> Promtail -> Grafana Loki -> Grafana
```

### otel-loki-tempo-mimir-grafana-no-promtail

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Tempo -> Grafana
```

#### Logs Flow
```text
Go App (file logs) -> OTel Collector (filelog receiver) -> Grafana Loki -> Grafana
```

### otel-promtail-loki-tempo-mimir-grafana

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Tempo -> Grafana
```

#### Logs Flow
```text
Go App (Zerolog) -> Promtail -> Grafana Loki -> Grafana
```

### otel-tempo-mimir-grafana

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Tempo -> Grafana
```

#### Logs Flow
```text
Not part of this branch's stack.
```

### otel-mimir-grafana

#### Metrics Flow
```text
Go App (OTel SDK) -> OTel Collector -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Not part of this branch's stack.
```

#### Logs Flow
```text
Not part of this branch's stack.
```

### prometheus-and-mimir

#### Metrics Flow
```text
Go App (Prometheus client) -> OTel Collector (Prometheus receiver) -> Grafana Mimir -> Grafana
```

#### Traces Flow
```text
Not part of this branch's stack.
```

#### Logs Flow
```text
Not part of this branch's stack.
```

## Configuration Files

- **docker-compose.yml**: Docker Compose configuration
- **otel-collector-config.yaml**: OpenTelemetry Collector configuration
- **configs/loki/loki-config.yml**: Loki configuration
- **mimir-config.yaml**: Mimir configuration
- **grafana/provisioning/**: Grafana provisioning files
- **hit-loop.sh**: Script to generate traffic

## Exploring in Grafana

Use Grafana's Explore feature to explore available telemetry:

1. Open Grafana (http://localhost:3000)
2. Click on the Explore icon in the left sidebar
3. Select a data source (Prometheus for metrics, Tempo for traces, Loki for logs)
4. Run queries and filter as needed

### Trace-to-Log Correlation

To follow a request through traces and logs:

1. Find a trace in Tempo
2. Click on any span
3. Click the "Logs for this span" button to see correlated logs
4. Alternatively, search logs with a trace ID: `{job="go-app"} |~ "trace_id=<trace-id>"`

## Exploring Metrics

Use Grafana's Explore feature to explore all available metrics:

1. Open Grafana (http://localhost:3000)
2. Click on the Explore icon in the left sidebar
3. Try queries like:
   - `{__name__=~".*request.*"}`
   - `rate(demo_requests_total[1m])`
   - `histogram_quantile(0.95, sum(rate(demo_request_duration_seconds_bucket[1m])) by (le))`

## Example PromQL Queries

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

## Why Two Approaches?

- **OpenTelemetry**: Newer standard that provides a unified approach for metrics, traces, and logs
- **Prometheus**: Well-established standard for metrics with a large ecosystem
- **Comparison**: Helps you decide which approach fits your needs better

## Switching Branches

To switch between the two implementations:

```bash
# Stop the current stack
docker-compose down

# Switch branch
git checkout prometheus-and-mimir  # or otel-mimir-grafana

# Start the stack again
docker-compose up -d
```

## Troubleshooting

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- OpenTelemetry community
- Prometheus community
- Grafana Labs for Mimir and Grafana

---

Made by [Sandeep Vishnu](https://github.com/sandeepkv93)
