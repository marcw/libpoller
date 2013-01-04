package main

import (
	"net/http"
	"time"
)

type Check struct {
	Url      string
	Key      string
	Interval time.Duration
	Header   http.Header
}

func NewCheck(url, key, duration string, headers map[string]string) (*Check, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}

	return &Check{Url: url, Key: key, Interval: d, Header: h}, nil
}

func (c *Check) Poll(client *http.Client) (int, time.Duration, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	req.Header = c.Header
	start := time.Now().UnixNano()
	resp, err := client.Do(req)
	end := time.Now().UnixNano()
	if err != nil {
		return 0, 0, err
	}

	return resp.StatusCode, time.Duration(end - start), nil
}
