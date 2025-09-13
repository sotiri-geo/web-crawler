package workers

import (
	"fmt"

	"github.com/sotiri-geo/web-crawler/internal/requests"
)

// Worker func - runs in its own go routine
func Worker(urlChannel chan string, resultChannel chan string, urlFetcher *requests.URLFetch) {
	// this iteration will occur until recive is closed
	for url := range urlChannel {
		// process url
		got, err := urlFetcher.FetchURL(url)
		if err != nil {
			// still push errors to the buffer queue
			resultChannel <- fmt.Sprintf("Error fetching %s: %v", url, err)
		} else {
			// Push to results channel
			resultChannel <- got.Content
		}
	}
}

func AggregateResults(urls []string, urlFetcher *requests.URLFetch, numWorkers int) []string {
	size := len(urls)
	urlChannel := make(chan string, size)
	resultChannel := make(chan string, size)
	go func() {
		for _, url := range urls {
			urlChannel <- url
		}
		close(urlChannel) // always remember to close channel
	}()

	// Set no. workers to process urls from urlChannel
	for i := 0; i < numWorkers; i++ {
		go Worker(urlChannel, resultChannel, urlFetcher)
	}
	// dont close the channel here: workers haven;t had time to write their results yet

	var results []string

	// push to results slice
	for i := 0; i < size; i++ {
		results = append(results, <-resultChannel)
	}

	return results
}
