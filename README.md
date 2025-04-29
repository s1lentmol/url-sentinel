# URL Sentinel

A lightweight, self-hosted service for monitoring the availability and uptime of your websites. URL Sentinel lets you register URLs, configure check intervals, and view real-time status and historical monitoring data via a simple REST API.

![Work Scheme](./work_scheme.png)

## Features

- **Register & Manage URLs**  
  Add, update or remove any number of URLs to monitor.

- **Customizable Check Intervals**  
  Define how often each URL should be probed (e.g. every 30s, 5m, 1h).

- **Background Monitoring**  
  Health checks run asynchronously in the background; no blocking of API requests.

- **Status & History**  
  • Retrieve the current status (HTTP code, response time) for each URL.  
  • View full history of past checks, with timestamps and details.

- **RESTful API**  
  Endpoints to manage URLs and fetch monitoring results in JSON.

- **Persistent Storage**  
  All monitoring data is stored in PostgreSQL for reliability and analysis.

---

## Technology Stack

- **Language:** Go  
- **Web Framework:** [Chi](https://github.com/go-chi/chi)  
- **Configuration:** [cleanenv](https://github.com/ilyakaznacheev/cleanenv)  
- **Structured Logging:** [slog](https://pkg.go.dev/log/slog)  
- **Database:** PostgreSQL

---

## Getting Started

1. **Clone the repo**  
   ```bash
   git clone https://github.com/s1lentmol/url-sentinel.git
   cd url-sentinel
   ```

2. **Configure environment**  
   Copy `.env.example` to `.env` and set your PostgreSQL connection string and other settings:
   ```ini
   DATABASE_URL=postgres://user:pass@localhost:5432/url_sentinel?sslmode=disable
   CHECK_INTERVAL_DEFAULT=60s
   PORT=8080
   ```

3. **Build & Run**  
   ```bash
   go build -o url-sentinel cmd/server/main.go
   ./url-sentinel
   ```

   The service will start on the port defined in your `.env` (default `8080`).

---

## REST API

### URLs

- **Create a new URL**  
  ```http
  POST /urls
  Content-Type: application/json

  {
    "address": "https://example.com",
    "interval": "30s"
  }
  ```
  **Response:** `201 Created`  
  ```json
  {
    "id": 1,
    "address": "https://example.com",
    "interval": "30s",
    "created_at": "2025-04-29T12:00:00Z"
  }
  ```

- **List all URLs**  
  ```http
  GET /urls
  ```
  **Response:** `200 OK`  
  ```json
  [
    {
      "id": 1,
      "address": "https://example.com",
      "interval": "30s",
      "created_at": "2025-04-29T12:00:00Z"
    },
    …
  ]
  ```

- **Get a single URL**  
  ```http
  GET /urls/{id}
  ```

- **Update check interval**  
  ```http
  PUT /urls/{id}
  Content-Type: application/json

  {
    "interval": "1m"
  }
  ```

- **Delete a URL**  
  ```http
  DELETE /urls/{id}
  ```

### Check History

- **List checks for a URL**  
  ```http
  GET /urls/{id}/checks
  ```
  **Response:** `200 OK`  
  ```json
  [
    {
      "timestamp": "2025-04-29T12:00:30Z",
      "status_code": 200,
      "response_time_ms": 123
    },
    …
  ]
  ```

- **Get latest check**  
  ```http
  GET /urls/{id}/checks/latest
  ```

---

## Database Schema

Two main tables:

- **urls**  
  - `id` (PK)  
  - `address` (TEXT)  
  - `interval` (INTERVAL)  
  - `created_at` (TIMESTAMP)

- **checks**  
  - `id` (PK)  
  - `url_id` (FK → urls.id)  
  - `timestamp` (TIMESTAMP)  
  - `status_code` (INT)  
  - `response_time_ms` (INT)

---

## How It Works

1. **Scheduler:** A background goroutine reads all registered URLs and their intervals.  
2. **Worker Pool:** Dedicated workers perform HTTP GET requests at the configured intervals.  
3. **Persistence:** Results (status code + latency) are inserted into PostgreSQL.  
4. **API:** Clients query the service at `/urls` endpoints to retrieve current status or historical data.
