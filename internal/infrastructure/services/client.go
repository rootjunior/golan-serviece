package services

import (
	"net/http"
)

type PostClient struct {
	URL    string
	Client *http.Client
}

func NewPostClient(url string) *PostClient {
	return &PostClient{
		URL:    url,
		Client: &http.Client{},
	}
}

func (c *PostClient) GetAll() {}
