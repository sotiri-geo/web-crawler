package internal_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sotiri-geo/web-crawler/internal"
)

type StubClient struct {
	response      string
	errorResponse error
}

func (c *StubClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(c.response)),
	}, nil
}

func TestFetchUrl(t *testing.T) {
	t.Run("fetches html", func(t *testing.T) {
		wantHtml := "<html><body>Hello World</body></html>"
		// Setup - basic server to return some canned response
		stubClient := &StubClient{response: wantHtml}
		urlFetcher := internal.URLFetcher{stubClient}

		got := urlFetcher.FetchURL("www.example.com")

		if got != wantHtml {
			t.Errorf("got %q, want %q", got, wantHtml)
		}
	})
}
