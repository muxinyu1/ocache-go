package impl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ocache/ocache"
)

type HttpClient struct {
	baseUrl string
	client  http.Client
}

func (c *HttpClient) Get(groupName string, key string) (ocache.Value, error) {
	path := fmt.Sprintf("%v/%v/%v", c.baseUrl, url.QueryEscape(groupName), url.QueryEscape(key))
	res, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return &ByteView{bytes}, nil
	}
	if res.StatusCode == http.StatusNotFound {
		return &ByteView{bytes}, fmt.Errorf("404") // TODO err中存在404并不是错误
	}
	return &ByteView{bytes}, nil
}
