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

## Prerequisites

You only need **Docker** installed on your machine to run this project.

### Install Docker

**Ubuntu / Debian / WSL:**
```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo usermod -aG docker $USER   # allow running docker without sudo
newgrp docker                   # apply group change in current shell
```

**macOS:** Download and install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

**Windows:** Install [Docker Desktop](https://www.docker.com/products/docker-desktop/) (WSL 2 backend recommended)

Verify:
```bash
docker --version         # e.g. Docker version 26.x.x
docker compose version   # e.g. Docker Compose version v2.x.x
```

---

## Quick Start

### Option A — Pull from GHCR (recommended)

No build needed, just pull and run:

```bash
docker run -d \
  --name uptime-monitor \
  -p 127.0.0.1:8080:8080 \
  -e CHECK_INTERVAL=60 \
  -v uptime_data:/data \
  ghcr.io/egayurcel990/go-uptime-monitor:latest
```

Open **http://localhost:8080** in your browser.

### Option B — Docker Compose

```bash
git clone https://github.com/egayurcel990/go-uptime-monitor
cd go-uptime-monitor
cp .env.example .env     # edit .env if you want to customize
docker compose up -d
```

Open **http://localhost:8080** in your browser.

### Option C — Build and run locally (requires Go 1.22+)

```bash
git clone https://github.com/egayurcel990/go-uptime-monitor
cd go-uptime-monitor
go build -o bin/monitor ./cmd/monitor
./bin/monitor
```

---

## Adding Your First Target

Once the app is running, add a URL to monitor:

```bash
curl -X POST http://localhost:8080/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name": "My Site", "url": "https://example.com", "interval": 60}'
```

The dashboard will automatically refresh and show your new target. You can also trigger an immediate check:

```bash
# Replace 1 with the target ID returned above
curl -X POST http://localhost:8080/api/v1/targets/1/check
```

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

## Input Validation

The API enforces the following rules on `POST /api/v1/targets`:

| Rule | Error response |
|------|----------------|
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

### Webhook Alerts

Set `WEBHOOK_URL` to a Slack or Discord incoming webhook to receive alerts when a target goes down. Alerts have a **5-minute cooldown per target** to prevent spam — you'll get one notification when it goes down, and another when it recovers.

**Slack:** Create an incoming webhook at https://api.slack.com/messaging/webhooks

**Discord:** Go to your server → Edit Channel → Integrations → Webhooks → New Webhook → Copy URL

---

## Stopping and Removing

```bash
# Stop the container
docker stop uptime-monitor

# Start it again (data is preserved)
docker start uptime-monitor

# Remove container and data completely
docker rm -f uptime-monitor
docker volume rm uptime_data
```

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
