<h1 align="center">go-uptime-monitor</h1>

<p align="center">
  Lightweight self-hosted URL uptime monitor — checks endpoints on a schedule, stores history, exposes results via HTTP API.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/SQLite-Database-003B57?style=flat-square&logo=sqlite&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/Prometheus-Metrics-E6522C?style=flat-square&logo=prometheus&logoColor=white" />
</p>

---

## What It Does

- Periodically checks a list of URLs (configurable interval)
- Records status code, latency, and up/down result to SQLite
- Exposes results via a simple HTTP API
- Exports `/metrics` endpoint for Prometheus scraping
- Sends alert to webhook (Slack / Discord) when a target goes down

---

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/targets` | List all monitored targets |
| `POST` | `/api/v1/targets` | Add a new target URL |
| `DELETE` | `/api/v1/targets/:id` | Remove a target |
| `GET` | `/api/v1/targets/:id/history` | Get check history for a target |
| `GET` | `/api/v1/status` | Overall uptime summary |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/healthz` | Health check |

---

## Quick Start

```bash
# Run with Docker
docker run -p 8080:8080 \
  -e CHECK_INTERVAL=60 \
  -v $(pwd)/data:/data \
  ghcr.io/egayurcel990/go-uptime-monitor:latest

# Or build and run locally
go build -o bin/monitor ./cmd/monitor
./bin/monitor
```

Add targets via API:
```bash
curl -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name": "My Site", "url": "https://example.com", "interval": 60}'
```

---

## Configuration

Via environment variables:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `/data/uptime.db` | SQLite database path |
| `CHECK_INTERVAL` | `60` | Default check interval (seconds) |
| `CHECK_TIMEOUT` | `10` | HTTP request timeout (seconds) |
| `WEBHOOK_URL` | — | Slack/Discord webhook for alerts |

---

## Project Structure

```
go-uptime-monitor/
├── cmd/
│   └── monitor/
│       └── main.go             # Entry point
├── internal/
│   ├── config/                 # Env-based config
│   ├── handler/                # HTTP handlers
│   ├── checker/                # URL check logic + scheduler
│   ├── repository/             # SQLite access layer
│   ├── model/                  # Domain types (Target, CheckResult)
│   ├── alert/                  # Webhook alert sender
│   └── metrics/                # Prometheus metrics
├── migrations/                 # SQL schema
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
make run      # Run locally
make test     # Run unit tests
make docker   # Build Docker image
make lint     # Run golangci-lint
```

---

## Prometheus Metrics

| Metric | Type | Description |
|---|---|---|
| `uptime_check_duration_seconds` | Histogram | HTTP check latency |
| `uptime_check_up` | Gauge | 1 = up, 0 = down per target |
| `uptime_checks_total` | Counter | Total checks performed |

Pair with Grafana for a full monitoring dashboard.

---

<p align="center">
  <i>Deployed via <a href="https://github.com/egayurcel990/ansible-server-bootstrap">ansible-server-bootstrap</a> · Universitas Brawijaya · 2025</i>
</p>
