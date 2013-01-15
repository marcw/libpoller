package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Contains a collection of Check
type ChecksList map[string]Check

// Default ChecksList
var Checks = make(ChecksList)

type Check struct {
	Url      *url.URL
	Addr     net.Addr
	Key      string
	Interval time.Duration
	Header   http.Header
}

// Represents the state of a Check after being polled
type CheckEvent struct {
	Check      *Check        // check
	Duration   time.Duration // total duration of check
	StatusCode int           // http status code, if any
	Time       time.Time     // time of check
	Timeout    bool
	Up         bool
}

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

func NewCheckEvent(check *Check) *CheckEvent {
	return &CheckEvent{Time: time.Now(), Check: check, Up: false, Timeout: false}
}

func (cl ChecksList) Add(c *Check) {
	cl[c.Key] = *c
}

func (cl ChecksList) AddFromJson(data []byte) error {
	checks := jsonChecks{}
	err := json.Unmarshal(data, &checks)
	if err != nil {
		return fmt.Errorf("There was an error reading your configuration file: %s", err)
	}

	for _, v := range checks {
		chk, err := NewCheck(v.Url, v.Key, v.Interval, v.Headers)
		if err != nil {
			return fmt.Errorf("Check configuration error: %s", err)
		}

		cl.Add(chk)
	}

	return nil
}

func (cl ChecksList) JSON() ([]byte, error) {
	checks := jsonChecks{}
	for _, v := range cl {
		header := make(map[string]string)
		for k, h := range v.Header {
			header[k] = h[0]
		}

		check := jsonCheck{Url: v.Url.String(), Key: v.Key, Interval: v.Interval.String(), Headers: header}
		checks = append(checks, check)
	}
	data, err := json.Marshal(checks)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (cl ChecksList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := cl.JSON()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Write(data)

		return
	}

	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)

			return
		}

		err = cl.AddFromJson(data)
		if err != nil {
			http.Error(w, err.Error(), 400)

			return
		}
	}

	// TODO: Need to implement PUT. But it means wiping the CheckList which can
	// cause problem when iterating on it in the main command.
}

func (cl ChecksList) Wipe() {
	for k, _ := range cl {
		delete(cl, k)
	}
}

func FromJson(data []byte) error {
	return Checks.AddFromJson(data)
}
