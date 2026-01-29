package main

import (
	"encoding/json"
	"io"
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

func (c *PostClient) GetPosts() ([]Post, error) {
	resp, err := http.Get(c.URL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}

	}(resp.Body)

	var posts []Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	return posts, err
}
