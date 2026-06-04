<h1 align="center">go-uptime-monitor</h1>

<p align="center">
  Lightweight self-hosted URL uptime monitor — checks endpoints on a schedule, stores history, exposes results via HTTP API and a built-in dashboard.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/SQLite-Database-003B57?style=flat-square&logo=sqlite&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/Prometheus-Metrics-E6522C?style=flat-square&logo=prometheus&logoColor=white" />
</p>

---

## What It Does

- Periodically checks a list of URLs (configurable interval, per-target or global)
- Records HTTP status code, latency, and up/down result to SQLite
- Serves a built-in web dashboard at `/`
- Exposes a REST API for managing targets and querying history
- Exports a `/metrics` endpoint for Prometheus scraping
- Sends webhook alerts (Slack / Discord) when a target goes down — with a **5-minute cooldown** to prevent alert spam
- Validates all input — rejects empty names, empty URLs, and non-http/https schemes with a clear error response

---

## Dashboard

![Dashboard](docs/dashboard.png)

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/targets` | List all monitored targets |
| `POST` | `/api/v1/targets` | Add a new target URL |
| `DELETE` | `/api/v1/targets/:id` | Remove a target |
| `GET` | `/api/v1/targets/:id/history` | Get check history (last 100) |
| `POST` | `/api/v1/targets/:id/check` | Trigger an immediate check |
| `GET` | `/api/v1/status` | Overall uptime summary (24h window) |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/healthz` | Health check |

---

## Quick Start

### Docker (recommended)

```bash
docker run -d \
  --name uptime-monitor \
  -p 127.0.0.1:8080:8080 \
  -e CHECK_INTERVAL=60 \
  -v $(pwd)/data:/data \
  ghcr.io/egayurcel990/go-uptime-monitor:latest
```

### Docker Compose

```bash
cp .env.example .env
# Edit .env if needed
docker compose up -d
```

### Build locally

```bash
go build -o bin/monitor ./cmd/monitor
./bin/monitor
```

Add a target:

```bash
curl -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name": "My Site", "url": "https://example.com", "interval": 60}'
```

---

## Input Validation

The API enforces the following rules on `POST /api/v1/targets`:

| Rule | Error response |
|------|---------------|
| `name` is empty | `{"error": "name is required"}` |
| `url` is empty | `{"error": "url is required"}` |
| `url` is not `http` or `https` | `{"error": "url must be a valid http or https URL"}` |

All validation errors return HTTP `400 Bad Request`.

---

## Configuration

All options are set via environment variables (or `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `/data/uptime.db` | SQLite database path |
| `CHECK_INTERVAL` | `60` | Default check interval in seconds |
| `CHECK_TIMEOUT` | `10` | HTTP request timeout in seconds |
| `WEBHOOK_URL` | — | Slack or Discord webhook URL for alerts |

---

## Project Structure

```
go-uptime-monitor/
├── cmd/
│   └── monitor/
│       └── main.go             # Entry point
├── internal/
│   ├── config/                 # Env-based configuration
│   ├── handler/                # HTTP handlers (Echo) + input validation
│   ├── checker/                # URL check logic + scheduler
│   ├── repository/             # SQLite data access layer
│   ├── model/                  # Domain types (Target, CheckResult, UptimeSummary)
│   ├── alert/                  # Webhook alert sender (5-min cooldown)
│   └── metrics/                # Prometheus metrics registration
├── web/
│   └── index.html              # Built-in dashboard
├── docs/                       # Screenshots
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── Makefile
└── README.md
```

---

## Makefile Commands

```bash
make build    # Compile binary to bin/monitor
make run      # Run locally with go run
make test     # Run unit tests
make docker   # Build Docker image
make lint     # Run golangci-lint
```

---

## Prometheus Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `uptime_check_duration_seconds` | Histogram | HTTP check latency per target |
| `uptime_check_up` | Gauge | 1 = up, 0 = down per target |
| `uptime_checks_total` | Counter | Total checks performed per target |

Pair with Grafana for a full monitoring dashboard.

---

## Deployment

This service is designed to be deployed via [ansible-server-bootstrap](https://github.com/egayurcel990/ansible-server-bootstrap), which sets up a hardened Ubuntu server with Nginx reverse proxy, UFW firewall, and Docker — then pulls and runs this image automatically.

---

<p align="center">
  <i>Deployed via <a href="https://github.com/egayurcel990/ansible-server-bootstrap">ansible-server-bootstrap</a> · Ega Yurcel Satriaji · 2025</i>
</p>
