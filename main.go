package main

import (
	"fmt"
	"io"
	"net/http"
)

type FetchResult struct {
	URL        string
	StatusCode int
	Size       int
	Error      error
}

func worker(id int, jobs <-chan string, results chan<- FetchResult) {
	for url := range jobs {
		resp, err := http.Get(url)
		if err != nil {
			results <- FetchResult{URL: url, Error: err}
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		results <- FetchResult{URL: resp.Request.URL.String(), StatusCode: resp.StatusCode, Size: len(body), Error: err}
	}
}

func main() {
	urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://uottawa.ca",
		"https://github.com",
		"https://httpbin.org/get",
	}

	numWorkers := 5

	jobs := make(chan string, len(urls))
	results := make(chan FetchResult, len(urls))

	//start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// send jobs
	for j := 1; j <= len(urls); j++ {
		jobs <- urls[j-1]
	}
	close(jobs)

	// collect results
	for i := 1; i <= len(urls); i++ {
		result := <-results

		if result.Error != nil {
			fmt.Println(result.URL, "| Error:", result.Error)
		}else{
			fmt.Println(result.URL, "| Status:", result.StatusCode, "| Size:", result.Size, "bytes")
		}
	}

	fmt.Println("\n Scraping complete!")
}
