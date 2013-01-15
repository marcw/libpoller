package check

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type Check struct {
	Url      *url.URL
	Addr     net.Addr
	Key      string
	Interval time.Duration
	Header   http.Header
}

// Represents the state of a Check after being polled
type Event struct {
	Check      *Check        // check
	Duration   time.Duration // total duration of check
	StatusCode int           // http status code, if any
	Time       time.Time     // time of check
	Up         bool
}

// Used for marshalling / unmarshalling
type jsonCheck struct {
	Url      string
	Key      string
	Interval string
	Headers  map[string]string
}

type jsonChecks []jsonCheck

func NewCheck(checkUrl, key, interval string, headers map[string]string) (*Check, error) {
	d, err := time.ParseDuration(interval)
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

	host := u.Host
	_, err = net.ResolveTCPAddr("tcp", u.Host)
	if err != nil {
		if u.Scheme == "http" {
			host = host + ":80"
		} else {
			host = host + ":443"
		}
	}

	a, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}

	return &Check{Url: u, Key: key, Interval: d, Header: h, Addr: a}, nil
}

func NewEvent(check *Check) *Event {
	return &Event{Time: time.Now(), Check: check, Up: false}
}
