package requests

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	ErrPageNotFound   = errors.New("could not find page")
	ErrInvalidURL     = errors.New("invalid url requested")
	ErrClientFetchURL = errors.New("client could not fetch url provided")
	ErrReadContent    = errors.New("cannot read body content from response")
)

// First build out a httpClient interface we can use as a dependency injection
// should mimic the interface found on http.Client
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type URLFetch struct {
	client       HTTPClient // dep injection allows us to create contracts with real implementations in integration tests
	urlValidator func(url string) bool
}

type Page struct {
	Content    string
	StatusCode int
}

func (u *URLFetch) FetchURL(url string) (*Page, error) {
	if !u.urlValidator(url) {
		return &Page{Content: "", StatusCode: http.StatusBadRequest}, ErrInvalidURL
	}
	response, err := u.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page %q: %w", url, ErrClientFetchURL)
	}
	defer response.Body.Close() // after successful resp
	if response.StatusCode == http.StatusNotFound {
		return &Page{Content: "", StatusCode: response.StatusCode}, ErrPageNotFound
	}
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Page{
		Content:    string(content),
		StatusCode: response.StatusCode,
	}, nil
}

// Set constructors
func NewURLFetch(client HTTPClient, urlValidator func(url string) bool) *URLFetch {
	return &URLFetch{client, urlValidator}
}

// Validator, separating this out from the FetchURL to decouple these two components
func URLValidator(url string) bool {
	pattern := `^www\.[A-Za-z0-9]+\.(com|co\.uk)$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}
