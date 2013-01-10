package check

import (
	"net/http"
	"net/url"
	"time"
)

type Check struct {
	Url      *url.URL
	Key      string
	Interval time.Duration
	Header   http.Header
}

func NewCheck(checkUrl, key, duration string, headers map[string]string) (*Check, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}

	u, err := url.Parse(checkUrl)
	if err != nil {
		return nil, err
	}

	return &Check{Url: u, Key: key, Interval: d, Header: h}, nil
}
