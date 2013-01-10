package config

import (
	"encoding/json"
	"fmt"
	"github.com/marcw/poller/backend"
	"github.com/marcw/poller/check"
	"os"
	"time"
)

const (
	DEFAULT_TIMEOUT    = "10s"
	DEFAULT_USER_AGENT = "Poller (https://github.com/marcw/poller)"
)

type Config struct {
	UserAgent string
	Url       string
	Timeout   time.Duration
	Backends  []backend.Backend
	Checks    []check.Check
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(data []byte) error {
	type ccheck struct {
		Url      string
		Key      string
		Interval string
		Headers  map[string]string
	}

	type configuration struct {
		Backends  []string
		UserAgent string
		Timeout   string
		Checks    []ccheck
	}

	config := &configuration{}

	err := json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("There was an error reading your configuration file: %s", err)
	}

	for _, v := range config.Backends {
		switch {
		case v == "statsd":
			statsd, err := backend.NewStatsdBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Statsd backend: %s", err)
			}
			c.Backends = append(c.Backends, statsd)

		case v == "stdout":
			stdout := backend.NewStdoutBackend()
			c.Backends = append(c.Backends, stdout)

		case v == "librato":
			librato, err := backend.NewLibratoBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Librato backend: %s", err)
			}
			c.Backends = append(c.Backends, librato)

		case v == "syslog":
			syslog, err := backend.NewSyslogBackend()
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
		chk, err := check.NewCheck(v.Url, v.Key, v.Interval, v.Headers)
		if err != nil {
			return fmt.Errorf("Check configuration error: %s", err)
		}
		chk.Header.Set("User-Agent", config.UserAgent)
		c.Checks = append(c.Checks, *chk)
	}

	c.UserAgent = config.UserAgent
	c.Timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		return fmt.Errorf("Invalid timeout value given: %s", err)
	}

	c.Url = os.Getenv("POLLER_URL")

	return nil
}

func (c *Config) CloseBackends() {
	for _, v := range c.Backends {
		v.Close()
	}
}
