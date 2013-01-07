package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	DEFAULT_TIMEOUT    = "10s"
	DEFAULT_USER_AGENT = "Poller (https://github.com/marcw/poller)"
)

type Configuration struct {
	UserAgent string
	Url       string
	Timeout   time.Duration
	Backends  []Backend
	Checks    []Check
}

func (c *Configuration) Load(data []byte) error {
	type check struct {
		Url      string
		Key      string
		Interval string
		Headers  map[string]string
	}

	type configuration struct {
		Backends  []string
		UserAgent string
		Timeout   string
		Checks    []check
	}

	config := &configuration{}

	err := json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("There was an error reading your configuration file: %s", err)
	}

	for _, v := range config.Backends {
		switch {
		case v == "statsd":
			statsd, err := NewStatsdBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Statsd backend: %s", err)
			}
			c.Backends = append(c.Backends, statsd)

		case v == "stdout":
			stdout := NewStdoutBackend()
			c.Backends = append(c.Backends, stdout)

		case v == "librato":
			librato, err := NewLibratoBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Librato backend: %s", err)
			}
			c.Backends = append(c.Backends, librato)

		case v == "syslog":
			syslog, err := NewSyslogBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the syslog backend: %s", err)
			}
			c.Backends = append(c.Backends, syslog)
		}
	}

	if config.Timeout == "" {
		config.Timeout = DEFAULT_TIMEOUT
	}

	if config.UserAgent == "" {
		config.UserAgent = DEFAULT_USER_AGENT
	}

	for _, v := range config.Checks {
		check, err := NewCheck(v.Url, v.Key, v.Interval, v.Headers)
		if err != nil {
			return fmt.Errorf("Check configuration error: %s", err)
		}
		check.Header.Set("User-Agent", config.UserAgent)
		c.Checks = append(c.Checks, *check)
	}

	c.UserAgent = config.UserAgent
	c.Timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		return fmt.Errorf("Invalid timeout value given: %s", err)
	}

	c.Url = os.Getenv("POLLER_URL")

	return nil
}

func (c *Configuration) CloseBackends() {
	for _, v := range c.Backends {
		v.Close()
	}
}
