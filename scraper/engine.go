package scraper

import (
	"log"
	"net/http"
	"time"
)

// Worker targets a single URL from the incoming queue channel and pipes the result out
func Worker(id int, urlQueue <-chan string, resultsOut chan<- ScrapeResult) {
	// Set up a shared corporate HTTP client with a built-in 5-second boundary safety timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// This loop keeps running, pulling tasks out of the channel until the channel is closed
	for url := range urlQueue {
		log.Printf("[Worker #%d] Jumping onto the internet to scan: %s", id, url)

		startTime := time.Now()
		resp, err := client.Get(url)

		// Handle network dropout or server downtime failures gracefully
		if err != nil {
			log.Printf("[Worker #%d] Network Error targeting %s: %v", id, url, err)
			resultsOut <- ScrapeResult{
				URL:        url,
				Title:      "Failed to reach domain",
				Price:      "0.00",
				ScrapedAt:  time.Now(),
				StatusCode: 500,
			}
			continue
		}

		// Always close the response streaming body to avoid memory leaks
		resp.Body.Close()

		// Simulate dynamic item rendering extraction logic
		var mockTitle string
		var mockPrice string
		
		switch id % 3 {
		case 0:
			mockTitle = "M4 Max MacBook Pro 14-inch Studio Display Edition"
			mockPrice = "1999.00"
		case 1:
			mockTitle = "Keychron Q1 Max Mechanical Keyboard (Custom Brown Switches)"
			mockPrice = "210.00"
		default:
			mockTitle = "Sony WH-1000XM5 Wireless Noise Canceling Headphones"
			mockPrice = "399.00"
		}

		log.Printf("[Worker #%d] ✅ Successfully harvested data from: %s (Time taken: %v)", id, url, time.Since(startTime))

		// Pipe the structured data model straight down the thread-safe results out channel
		resultsOut <- ScrapeResult{
			URL:        url,
			Title:      mockTitle,
			Price:      mockPrice,
			ScrapedAt:  time.Now(),
			StatusCode: resp.StatusCode,
		}
	}
}
