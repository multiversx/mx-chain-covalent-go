package api

import (
	"net/http"
	"time"
)

type httpClient struct {
	client *http.Client
}

// NewDefaultHttpClient will create a default http client which can be queried to GET url responses
func NewDefaultHttpClient(requestTimeoutSec uint64) *httpClient {
	client := &http.Client{}
	client.Timeout = time.Duration(requestTimeoutSec) * time.Second

	return &httpClient{
		client: client,
	}
}

// Get will return an resp *http.Response from the provided url
func (hc *httpClient) Get(url string) (resp *http.Response, err error) {
	return hc.client.Get(url)
}
