package requests

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrPageNotFound   = errors.New("could not find page")
	ErrInvalidURL     = errors.New("invalid urlPage requested")
	ErrClientFetchURL = errors.New("client could not fetch urlPage provided")
	ErrReadContent    = errors.New("cannot read body content from response")
)

// First build out a httpClient interface we can use as a dependency injection
// should mimic the interface found on http.Client
type HTTPClient interface {
	Get(urlPage string) (*http.Response, error)
}

type URLFetch struct {
	client       HTTPClient // dep injection allows us to create contracts with real implementations in integration tests
	urlValidator func(urlPage string) bool
}

type Page struct {
	Content    string
	StatusCode int
}

func (u *URLFetch) FetchURL(urlPage string) (*Page, error) {
	if !u.urlValidator(urlPage) {
		return &Page{Content: "", StatusCode: http.StatusBadRequest}, ErrInvalidURL
	}
	response, err := u.client.Get(urlPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page %q: %w", urlPage, ErrClientFetchURL)
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
func NewURLFetch(client HTTPClient, urlValidator func(urlPage string) bool) *URLFetch {
	return &URLFetch{client, urlValidator}
}

// Validator, separating this out from the FetchURL to decouple these two components
func URLValidator(urlPage string) bool {
	parsedURL, err := url.Parse(urlPage)

	if err != nil {
		return false
	}

	// Check valid schema
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return true
}
