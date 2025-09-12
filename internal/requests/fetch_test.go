package requests_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sotiri-geo/web-crawler/internal/requests"
)

// Mock needs to implement the HTTPClient interface
type MockClient struct {
	content    string // canned html content responses
	statusCode int
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.content)),
	}, nil
}

func TestFetchURL(t *testing.T) {
	t.Run("fetches html content from page", func(t *testing.T) {
		url := "www.example.com"
		want := "<html><body>Hello World</body></html>"
		mockClient := &MockClient{content: want, statusCode: http.StatusOK}
		urlFetch := requests.NewURLFetch(mockClient)

		got, err := urlFetch.FetchURL(url)

		if err != nil {
			t.Fatalf("should not fail: found %v", err)
		}

		if got.StatusCode != http.StatusOK {
			t.Fatalf("incorrect status code: got %d, want %d", got.StatusCode, http.StatusOK)
		}

		if got.Content != want {
			t.Errorf("incorrect content returned: got %q, want %q", got, want)
		}
	})

	t.Run("page not found returns 404 with empty content", func(t *testing.T) {
		url := "www.notfound.com"
		want := ""
		mockClient := &MockClient{content: want, statusCode: http.StatusNotFound}
		urlFetch := requests.NewURLFetch(mockClient)

		got, err := urlFetch.FetchURL(url)

		// We need to return a custom client side error type when not found
		if !errors.Is(err, requests.ErrPageNotFound) {
			t.Fatalf("incorrect error returned: got %v, want %v", err, requests.ErrPageNotFound)
		}

		if got.Content != want {
			t.Errorf("incorrect html content found: got %q, want %q", got.Content, want)
		}

	})
}

func TestURLValidator(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{name: "valid url with .com", url: "www.example.com", want: true},
		{name: "valid url with .co.uk", url: "www.example.co.uk", want: true},
		{name: "missing www.", url: "example.com", want: false},
		{name: "missing domain", url: "www.example", want: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if isValid := requests.URLValidator(tt.url); isValid != tt.want {
				t.Errorf("incorrectly marked %q as valid = %v", tt.url, isValid)
			}
		})
	}
}
