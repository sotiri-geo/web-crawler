package workers_test

import (
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/sotiri-geo/web-crawler/internal/requests"
	"github.com/sotiri-geo/web-crawler/internal/workers"
)

type MockClient struct {
	content string
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

	t.Run("processes single url from url channel", func(t *testing.T) {
		// setup: Single worker and url to validate process
		urlChannel := make(chan string, 1)
		resultChannel := make(chan string, 1)
		want := "<html><body>Hello World</body></html>"
		mockClient := &MockClient{content: "<html><body>Hello World</body></html>"}
		urlFetcher := requests.NewURLFetch(mockClient, func(urlPage string) bool { return true })

		// execute: background worker
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
}
