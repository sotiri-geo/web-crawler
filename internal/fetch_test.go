package internal_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sotiri-geo/web-crawler/internal"
)

// Testing an interface
type StubClient struct {
	contentResponse string
	errorResponse   error
}

func (c *StubClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(c.contentResponse)),
	}, nil
}

func NewStubClient(contentResponse string, errorResponse error) *StubClient {
	return &StubClient{contentResponse, errorResponse}
}

func TestFetchUrl(t *testing.T) {
	// Start by refactoring into table tests with setup

	cases := []struct {
		name              string
		url               string
		wantHtmlContent   string
		wantErrorResponse error
		wantStatusCode    int
	}{
		{
			name:              "fetches html content",
			url:               "www.example.com",
			wantHtmlContent:   "<html><body>Hello World</body></html>",
			wantErrorResponse: nil,
			wantStatusCode:    http.StatusOK,
		},
	}

	for _, tt := range cases {
		// setup
		client := NewStubClient(tt.wantHtmlContent, tt.wantErrorResponse)
		urlFetcher := internal.NewURLFetcher(client)

		// execute
		got, err := urlFetcher.FetchURL(tt.url)

		if tt.wantErrorResponse == nil {
			assertNoErr(t, err)
		}

		// assert
		assertStatusCode(t, got.StatusCode, tt.wantStatusCode)
		assertHtmlContent(t, got.Content, tt.wantHtmlContent)
		assertURL(t, got.URL, tt.url)

	}

	t.Run("returns empty string on error", func(t *testing.T) {

	})

	t.Run("handles 404 response", func(t *testing.T) {})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("failed status code: got %d, want %d", got, want)
	}
}

func assertHtmlContent(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("html content: got %s, want %s", got, want)
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("should not error: %v", err)
	}
}

func assertURL(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("incorrect URL call: got %q, want %q", got, want)
	}
}
