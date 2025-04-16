# Go Metrics Instrumentation and Visualization Demo 📊 📈

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## 🔍 Overview

This repository demonstrates two different approaches to instrumenting Go applications with metrics, showing both the modern OpenTelemetry approach and the traditional Prometheus client approach.

## 🌟 Features

- ✅ Complete working examples of metrics instrumentation
- ✅ Docker Compose setup for quick deployment
- ✅ Pre-configured Grafana dashboards
- ✅ Request count and latency histogram metrics
- ✅ Side-by-side comparison of different instrumentation approaches

## 🔄 Branches

### 🚀 otel-mimir-grafana (default)

**Modern OpenTelemetry Approach**

This branch demonstrates the modern approach to metrics using the OpenTelemetry standard:

- **Instrumentation**: Uses OpenTelemetry SDK for Go
- **Collection**: OpenTelemetry Collector
- **Storage**: Grafana Mimir (Prometheus-compatible)
- **Visualization**: Grafana dashboards

**Data Flow**: Go App → OTel SDK → OTel Collector → Mimir → Grafana

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

   ```
   http://localhost:3000
   ```

   Default credentials:

   ```
   Username: admin
   Password: admin
   ```

4. Generate some test traffic:

   ```bash
   ./hit-loop.sh
   ```

   Or create and use this script:

   ```bash
   #!/bin/bash

   # Create hit-loop.sh file
   cat > hit-loop.sh << 'EOF'
   #!/bin/bash

   # Number of requests to send
   REQUESTS=${1:-100}
   # Delay between requests in seconds
   DELAY=${2:-0.1}
   # URL to hit
   URL=${3:-"http://localhost:8080/metrics"}

   echo "Sending $REQUESTS requests to $URL with ${DELAY}s delay"

   for i in $(seq 1 $REQUESTS); do
     echo "Request $i/$REQUESTS"
     curl -s "$URL" > /dev/null
     sleep $DELAY
   done

   echo "Done!"
   EOF

   # Make it executable
   chmod +x hit-loop.sh

   # Run it
   ./hit-loop.sh 200 0.2
   ```

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

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. Commit your changes:
   ```bash
   git commit -m 'Add some amazing feature'
   ```
4. Push to the branch:
   ```bash
   git push origin feature/amazing-feature
   ```
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- OpenTelemetry community
- Prometheus community
- Grafana Labs for Mimir and Grafana

---

Made with ❤️ by [Sandeep Vishnu](https://github.com/sandeepkv93)
