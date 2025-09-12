package internal

import (
	"fmt"
	"io"
)

type URLFetcher struct {
	Client HTTPClient
}

type FetchResult struct {
	URL        string
	Content    string // html
	StatusCode int
	Error      error
}

func (u *URLFetcher) FetchURL(url string) (*FetchResult, error) {
	resp, err := u.Client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("client failed to fetch url: %v", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to ready body content: %v", err)
	}
	return &FetchResult{
		URL:        url,
		Content:    string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

func NewURLFetcher(client HTTPClient) *URLFetcher {
	return &URLFetcher{client}
}
