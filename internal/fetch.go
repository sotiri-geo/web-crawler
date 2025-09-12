package internal

import "io"

type URLFetcher struct {
	Client HTTPClient
}

type URLResponse struct {
	Html       string
	StatusCode int
}

func (u *URLFetcher) FetchURL(url string) *URLResponse {
	resp, _ := u.Client.Get(url)

	body, _ := io.ReadAll(resp.Body)
	return &URLResponse{
		Html:       string(body),
		StatusCode: resp.StatusCode,
	}
}
