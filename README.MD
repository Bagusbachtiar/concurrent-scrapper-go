# Go High-Performance Concurrent Web Scraper REST API

A multi-threaded, concurrent web scraping microservice and data ingestion API built natively in Go. It exposes `POST /scrape` and `GET /products` HTTP endpoints via Go's `net/http` library, utilizing worker goroutines, buffered channels, `goquery` for DOM parsing, a Neon Serverless Postgres cluster for atomic data persistence, and a multi-stage Alpine Docker container.

---

### 🛠️ System Architecture & Features

*   **RESTful HTTP Microservice (`net/http`):** Listens on port `8080`, handling JSON payloads (`ScrapeRequest`) and returning structured metrics (`ScrapeResponse`).
*   **Dynamic Concurrent Worker Pool:** Uses background `goroutines` and typed channels (`urlQueue`, `resultsOut`) for race-condition-free parallel processing.
*   **Live HTML DOM Parsing:** Extracts semantic metadata (such as `<title>` tags) and custom `User-Agent` headers using `goquery`.
*   **Atomic Cloud Persistence:** Connects to Neon Serverless Postgres via `pgxpool` (`default_query_exec_mode=simple_protocol`) to execute conflict-safe `UPSERT` queries.
*   **Multi-Stage Docker Containerization:** Uses a two-stage `Dockerfile` (`golang:alpine` builder + `alpine:latest` runtime) to compile a stripped static Linux binary (~15MB–20MB).

---

### 🐳 Running with Docker

Ensure **Docker Desktop** is running and your `.env` file contains a valid `DATABASE_URL`:

```bash
# 1. Build the lightweight production image
docker build -t go-scraper-api .

# 2. Run the container with runtime environment secrets
docker run --rm -p 8080:8080 --env-file .env go-scraper-api
```

---

### 🔌 API Endpoint Usage

#### `POST /scrape`
Accepts a JSON array of target URLs, crawls them concurrently, saves records to Neon Postgres, and returns execution metrics.

**Request Payload:**
```json
{
  "urls": [
    "https://go.dev",
    "https://github.com",
    "https://ycombinator.com"
  ]
}
```

**Response Payload (`200 OK`):**
```json
{
  "total_scraped": 2,
  "time_taken": "440.716292ms",
  "results": [
    {
      "url": "https://go.dev",
      "title": "The Go Programming Language",
      "price": "149.99",
      "scraped_at": "2026-11-13T12:06:40Z",
      "status_code": 200
    }
  ]
}
```

#### `GET /products`
Fetches the 20 most recently scraped records from the Neon Postgres `scraped_products` table.

**Response Payload (`200 OK`):**
```json
[
  {
    "url": "https://github.com",
    "title": "GitHub · Change is constant. GitHub keeps you ahead.",
    "price": "149.99",
    "scraped_at": "2026-11-12T17:22:19Z",
    "status_code": 200
  }
]
```

---

### 📊 Verified Live Server Logs

Containerized execution logs demonstrating concurrent URL scraping, database commit operations, and response timing can be found in your local server output, completing a 2-URL live batch inside Alpine Linux in **440ms**.
