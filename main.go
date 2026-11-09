package main 

import (
	"concurrent-scraper/scraper"
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main (){
	log.Println("🕷️ Concurrent Web Scraper Engine Initiating...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found, reading from system hardware environtment parameters.")
	}

	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == ""{
		log.Fatalf("Critical Error: DATABASE_URL is completely missing from environtment allocation")
	}

	//Connect to neon cluster
	dbPool, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to allocate pool link resource: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil{
		log.Fatalf("Database cloud server handshake failed: %v", err)
	}
	log.Println("Successfully connected to Database")

	targetURLs := []string{
		"https://httpbin.org",
		"https://httpbin.org",
		"https://httpbin.org",
		"https://httpbin.org",
		"https://httpbin.org",
	}

	//1. Initialize our thread safe channel pipelines (paralon pipings)
	urlQueue := make(chan string, len (targetURLs))
	resultsOut := make(chan scraper.ScrapeResult, len(targetURLs))

	//2. Feed all tasks into our inbound URL Pipleine channel
	for _, url := range targetURLs{
		urlQueue <- url
	}

	close (urlQueue)// closing tells the workers no more items are coming down the wire 

	totalWorkers := 3
	log.Printf("🚀 Spawning a thread-safe worker pool of %d concurrent goroutines...", totalWorkers)

	startTime := time.Now()

	//3. Fire up our background workers concurrently using the 'go' keyword 
	for i :=1; i <= totalWorkers; i++ {
		go scraper.Worker (i, dbPool, urlQueue, resultsOut)
	}

	//4. Stand at the bottom of the results out piepline and harvest data packets as they drop out
	 log.Println("📥 Collecting extracted metrics data packets streaming out of channels...")
	for i := 0; i < len(targetURLs); i++ {
		result := <-resultsOut
		fmt.Printf("\n📦 [Crawl Result Captured]\n - URL: %s\n - Title: %s\n - Price: $%s\n - Code: %d\n",
				result.URL, result.Title, result.Price, result.StatusCode)
	}

	log.Printf("\n🏁 Job Finished. Succefully crawled %d high-latency pages in: %v", len(targetURLs), time.Since(startTime))
}	