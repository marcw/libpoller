package main

import (
    "time"
    "net/http"
)

type Check struct {
    Url      string
    Key      string
    Interval time.Duration
}

func NewCheck(url, key, duration string) (*Check, error) {
    d, err := time.ParseDuration(duration)
    if err != nil {
        return nil, err
    }

    return &Check{Url: url, Key: key, Interval: d}, nil
}

func (c *Check) Poll(client *http.Client) (int, time.Duration, error) {
    start := time.Now().UnixNano()
    resp, err := client.Get(c.Url)
    end := time.Now().UnixNano()
    if err != nil {
        return 0, 0, err
    }

    return resp.StatusCode, time.Duration(end - start), nil
}

