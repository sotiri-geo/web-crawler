package workers

import (
	"fmt"

	"github.com/sotiri-geo/web-crawler/internal/requests"
)

func Worker(urlChannel chan string, resultChannel chan string, urlFetcher *requests.URLFetch) {
	for url := range urlChannel {
		// process url
		got, err := urlFetcher.FetchURL(url)
		if err != nil {
			// still push errors to the buffer queue
			resultChannel <- fmt.Sprintf("Error fetching %s: %v", url, err)
		}
		// Push to results channel
		resultChannel <- got.Content
	}
}
