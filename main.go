package main 

import (
	"concurrent-scraper/scraper"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("🕷️ Concurrent Web Scraper API Initializing...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found, reading from system environment.")
	}

	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		log.Fatalf("Critical Error: DATABASE_URL is missing from environment.")
	}

	dbPool, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("✅ Successfully connected to Neon serverless Postgres Cluster!")

	// Register our POST /scrape API endpoint
	http.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {
		// 1. Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
			return
		}

		// 2. Decode incoming JSON body into our ScrapeRequest struct
		var reqPayload scraper.ScrapeRequest
		if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if len(reqPayload.URLs) == 0 {
			http.Error(w, "Please provide at least one URL in the 'urls' array", http.StatusBadRequest)
			return
		}

		log.Printf("📥 [API Request] Received batch of %d URLs to scrape concurrently...", len(reqPayload.URLs))
		startTime := time.Now()

		// 3. Initialize channels dynamically based on incoming request length
		urlQueue := make(chan string, len(reqPayload.URLs))
		resultsOut := make(chan scraper.ScrapeResult, len(reqPayload.URLs))

		for _, url := range reqPayload.URLs {
			urlQueue <- url
		}
		close(urlQueue)

		// 4. Spawn 3 concurrent workers using your existing engine
		totalWorkers := 3
		for i := 1; i <= totalWorkers; i++ {
			go scraper.Worker(i, dbPool, urlQueue, resultsOut)
		}

		// 5. Gather all results from the output channel into a slice
		var collectedResults []scraper.ScrapeResult
		for i := 0; i < len(reqPayload.URLs); i++ {
			res := <-resultsOut
			collectedResults = append(collectedResults, res)
		}

		// 6. Build and send the final JSON response
		response := scraper.ScrapeResponse{
			TotalScraped: len(collectedResults),
			TimeTaken:    time.Since(startTime).String(),
			Results:      collectedResults,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		log.Printf("🏁 [API Response] Sent %d scraped records back to client in %v", len(collectedResults), time.Since(startTime))
	})

	// Start listening on port 8080
	port := ":8080"
	log.Printf("🚀 Scraper API Server is live and listening on http://localhost%s/scrape", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
