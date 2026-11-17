package main

import "fmt"

func main() {
	logs := []string{"INFO", "ERROR", "INFO", "WARN", "ERROR", "ERROR"}

	// TODO 1: Create an empty map named 'counts' with string keys and int values
	counts["string"] = 6
	// TODO 2: Loop through the 'logs' slice using 'range' and increment the count for each status
	for logs, total := range counts {
		gmt.Println(logs, ":", total)
	}
	// TODO 3: Print out the total count for "ERROR"
}
