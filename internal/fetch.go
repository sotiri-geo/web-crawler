package internal

import "io"

type URLFetcher struct {
	Client HTTPClient
}

func (u *URLFetcher) FetchURL(url string) string {
	resp, _ := u.Client.Get(url)

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
