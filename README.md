# URL Sentinel

A service for monitoring URL availability with automatic health checks and result history.

## Description

URL Sentinel is a REST API service that periodically checks the availability of added URLs and stores check results. Built with Clean Architecture principles and production-ready.

## Tech Stack

- **Go 1.24.1** — core language
- **PostgreSQL 15** — data storage
- **Chi Router** — HTTP routing
- **Docker & Docker Compose** — containerization
- **Clean Architecture** — architectural approach
- **lib/pq** — PostgreSQL driver
- **cleanenv** — configuration management
- **slog** — structured logging

## Quick Start

```bash
# Clone repository
git clone https://github.com/s1lentmol/url-sentinel.git
cd url-sentinel

# Run with Docker
docker-compose up -d

# Or locally
make run
```

## API

- `POST /url` — add URL for monitoring
- `GET /url/{id}` — get URL information
- `GET /url/list` — list all URLs
- `DELETE /url/{id}` — delete URL
- `GET /url/{id}/history` — URL check history

