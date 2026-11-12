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

//ScrapeRequest models the incoming JSON payload sent to our POST /scrape endpoint
type ScrapeRequest struct {
	URLs []string `json:"urls"`
}

//ScrapeReponse models the structured JSON reply sent back to the caller 
type ScrapeResponse struct {
	TotalScraped	int				`json:"total_scraped"`
	TimeTaken		string			`json:"time_taken"`
	Results			[]ScrapeResult	`json:"results"`
}