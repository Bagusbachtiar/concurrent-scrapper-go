package scraper 
 
import "time"

//ScrapeResult models the data layout harvested from our target pagess
type ScrapeResult struct {
	URL 	string	`json:"url"`
	Title 	string 	`json:"title"`
	Price 	string 	`jon:"price"`
	ScrapedAt	time.Time `json:"scraped_at"`
	StatusCode 	int 	`json:"status_code"`
}