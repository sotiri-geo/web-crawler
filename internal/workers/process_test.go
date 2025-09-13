package workers_test

import (
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/sotiri-geo/web-crawler/internal/requests"
	"github.com/sotiri-geo/web-crawler/internal/workers"
)

type MockClient struct {
	content string
	delay   time.Duration
}

func (m *MockClient) Get(urlPage string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(m.content)),
	}, nil
}

func TestWorkerPool(t *testing.T) {
	t.Run("send urls to workers", func(t *testing.T) {
		// This essentially tests the plumbing of the worker pool. I.e. channels PUSH, POP operations
		urlChannel := make(chan string, 3) // creates a buffered channel to hold up to 3 strings. Sending won't block until buffer is full

		// When: We send URLS to the channel
		urls := []string{"http://site1.com", "http://site2.com"}

		// Single producer: pushes onto buffer queue (which is channel)
		go func() {
			for _, url := range urls {
				urlChannel <- url // send url to channel (buffer)
			}
			close(urlChannel) // close after completion of adding all urls to channel - "no more urls coming"
		}()

		// Read off the channel to make sure all urls are there for consumers to read from
		// Then: workers should recieve all URLS (i.e. the goroutines which read off the channel)
		var recievedURLs []string

		for url := range urlChannel {
			recievedURLs = append(recievedURLs, url)
		}

		// Check all the urls are in the correct order
		// Asserts correct order of URLS to be popped from queue (i.e. channel)
		if !slices.Equal(urls, recievedURLs) {
			t.Errorf("URLs not added to buffer channel: got %v, want %v", recievedURLs, urls)
		}
	})

	t.Run("processes url from channel", func(t *testing.T) {
		// setup: Single worker and url to validate process
		urlChannel := make(chan string, 1)
		resultChannel := make(chan string, 1)
		want := "<html><body>Hello World</body></html>"
		mockClient := &MockClient{content: "<html><body>Hello World</body></html>"}
		urlFetcher := requests.NewURLFetch(mockClient, func(urlPage string) bool { return true })

		// execute: background worker: it will process result and shuve html content into the resultChannel
		go workers.Worker(urlChannel, resultChannel, urlFetcher)

		// Send to channel a new url to be processed by worker
		urlChannel <- "http://www.example.com"
		close(urlChannel) // close after sending all urls to channel

		// Assert - read result from the resultChannel
		got := <-resultChannel

		if got != want {
			t.Errorf("incorrect result fetched from channel: got content %q, want %q", got, want)
		}
	})

	t.Run("process multiple urls by 2 background workers", func(t *testing.T) {
		// setup
		numWorkers := 2
		urls := []string{
			"http://www.example1.com",
			"http://www.example2.com",
			"http://www.example3.com",
			"http://www.example4.com",
		}

		htmlContent := "<html><body>Hello World</body></html>"
		wantResults := []string{
			htmlContent,
			htmlContent,
			htmlContent,
			htmlContent,
		}

		urlChannel := make(chan string, len(urls))
		resultChannel := make(chan string, len(urls))

		mockClient := &MockClient{content: htmlContent, delay: (time.Millisecond * 100)} // simulate network delays
		urlFetcher := requests.NewURLFetch(mockClient, func(url string) bool { return true })

		// execute: add urls to channel and create 2 workers
		// Simulate work
		go func() {
			for _, url := range urls {
				urlChannel <- url
			}
			close(urlChannel) // stop reading off channel after urls have been added
		}()

		// start workers
		for i := 0; i < numWorkers; i++ {
			go workers.Worker(urlChannel, resultChannel, urlFetcher)
		}

		var gotResults []string
		start := time.Now()
		for i := 0; i < len(urls); i++ {
			// add to results slice from channel where workers dump results
			gotResults = append(gotResults, <-resultChannel)
		}

		end := time.Now()
		totalElapsedTime := end.Sub(start)

		// Assert: verify concurrency by duration taken and validated results
		if totalElapsedTime >= time.Millisecond*300 {
			t.Errorf("took too long to run with simulated network delays: got %v and want under %v", totalElapsedTime, time.Millisecond*300)
		}

		if !slices.Equal(gotResults, wantResults) {
			t.Errorf("incorrect content found: got %v, want %v", gotResults, wantResults)
		}
	})
}
