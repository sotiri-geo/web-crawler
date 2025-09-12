package requests

import (
	"io"
	"net/http"
)

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
