package requests

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

var ErrPageNotFound = errors.New("could not find page")

// First build out a httpClient interface we can use as a dependency injection
// should mimic the interface found on http.Client
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type URLFetch struct {
	client HTTPClient // dep injection allows us to create contracts with real implementations in integration tests
}

type Page struct {
	Content    string
	StatusCode int
}

func (u *URLFetch) FetchURL(url string) (*Page, error) {
	response, _ := u.client.Get(url)
	if response.StatusCode == http.StatusNotFound {
		return &Page{Content: "", StatusCode: response.StatusCode}, ErrPageNotFound
	}
	content, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	return &Page{
		Content:    string(content),
		StatusCode: response.StatusCode,
	}, nil
}

// Set constructors
func NewURLFetch(client HTTPClient) *URLFetch {
	return &URLFetch{client}
}

// Validator, separating this out from the FetchURL to decouple these two components
func URLValidator(url string) bool {
	pattern := `^www\.[A-Za-z0-9]+\.(com|co\.uk)$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}
