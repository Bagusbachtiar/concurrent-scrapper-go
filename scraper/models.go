package craper 
 
import "time"

//ScrapeResult models the data layout harvested from our target pagess
type ScrapeResult struct {
	URL 	string	`json:"url"`
	Title 	string 	`json:"title"`
	Price 	string 	`jon:"price"`
	Scrapedat	string `json:"scrapeda_at"`
	StatusCode 	int 	`json:"status_code"`
}