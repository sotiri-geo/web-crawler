package internal

import (
	"fmt"
	"io"
)

type URLFetcher struct {
	Client HTTPClient
}

type URLResponse struct {
	HtmlContent string
	StatusCode  int
}

func (u *URLFetcher) FetchURL(url string) (*URLResponse, error) {
	resp, err := u.Client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("client failed to fetch url: %v", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to ready body content: %v", err)
	}
	return &URLResponse{
		HtmlContent: string(body),
		StatusCode:  resp.StatusCode,
	}, nil
}
