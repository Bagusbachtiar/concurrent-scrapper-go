# Go High-Performance Concurrent Web Scraper & Data Ingestion Pipeline

A multi-threaded, concurrent web scraping engine and data ingestion pipeline built natively in Go. This service utilizes buffered channels and lightweight worker goroutines to crawl live web endpoints simultaneously, parse raw HTML DOM elements, and stream structured metrics directly into a cloud database cluster.

---

### 🛠️ System Architecture & Features

*   **Concurrent Worker Pool:** Spawns isolated background `goroutines` that pull target URLs from a shared queue channel, allowing multiple network requests to execute in parallel without blocking the main thread.
*   **Thread-Safe Channel Pipelines:** Uses typed Go channels (`urlQueue` and `resultsOut`) to safely coordinate task distribution and result collection with zero race conditions.
*   **Live HTML DOM Parsing:** Integrates `goquery` to load raw HTTP response streams and extract semantic HTML metadata (`<title>` tags) while injecting custom browser `User-Agent` headers.
*   **Atomic Cloud Persistence:** Connects to a **Neon Serverless Postgres** cluster via `pgxpool` using simple protocol mode (`default_query_exec_mode=simple_protocol`) to execute conflict-safe `UPSERT` queries (`ON CONFLICT DO UPDATE`).
*   **Resilient Network Boundaries:** Enforces strict HTTP client timeouts (`7s`) and immediate response body closures (`resp.Body.Close()`) to prevent memory leaks and hanging threads.

---

### 📊 Verified Live Execution Logs

During live test runs crawling production web endpoints, the 3-worker pool harvested live HTML titles and committed all records to Neon Postgres in **1.23 seconds**:

```text
🕷️ Concurrent Web Scraper Engine Initiating...
Successfully connected to Database
🚀 Spawning a thread-safe worker pool of 3 concurrent goroutines...
📥 Collecting extracted metrics data packets streaming out of channels...
[Worker #1] 🕷️ Initiating live extraction sequence on: https://go.dev
[Worker #2] 🕷️ Initiating live extraction sequence on: https://github.com
[Worker #3] 🕷️ Initiating live extraction sequence on: https://wikipedia.org
[Worker #2] 💾 Data Permanently committed to Neon Database for : https://github.com
[Worker #2] ✅ Successfully Finished live task parsing for https://github.com in 326.7134ms
[Worker #3] 💾 Data Permanently committed to Neon Database for : https://wikipedia.org
[Worker #3] ✅ Successfully Finished live task parsing for https://wikipedia.org in 377.1769ms
[Worker #1] 💾 Data Permanently committed to Neon Database for : https://go.dev
[Worker #1] ✅ Successfully Finished live task parsing for https://go.dev in 576.0659ms
[Worker #2] 💾 Data Permanently committed to Neon Database for : https://httpbin.org
[Worker #2] ✅ Successfully Finished live task parsing for https://httpbin.org in 907.9938ms

🏁 Job Finished. Successfully crawled 4 high-latency pages in: 1.2362106s
```
