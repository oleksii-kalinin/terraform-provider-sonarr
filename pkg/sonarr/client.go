package sonarr

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	ApiKey     string
	HttpClient *http.Client
}

func NewClient(url, key string) *Client {
	timeout := 10 * time.Second // Hardcode for now

	return &Client{
		BaseURL: url,
		ApiKey:  key,
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) doRequest(r *http.Request) (*http.Response, error) {
	r.Header.Add("X-Api-Key", c.ApiKey)
	r.Header.Add("Content-Type", "application/json")

	res, err := c.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	return res, err
}
