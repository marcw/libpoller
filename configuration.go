package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Configuration struct {
	UserAgent string
	Url       string
	Timeout   time.Duration
	Backends  []Backend
	Checks    []Check
}

func (c *Configuration) Load(data []byte) {
	type check struct {
		Url      string
		Key      string
		Interval string
		Headers  map[string]string
	}

	type configuration struct {
		UserAgent string
		Timeout   string
		Backends  []string
		Checks    []check
	}

	config := &configuration{}

	err := json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("There was an error reading your configuration file:", err)
	}

	for _, v := range config.Backends {
		switch {
		case v == "statsd":
			statsd, err := NewStatsdBackend()
			if err != nil {
				log.Fatalln("Impossible to instanciate the Statsd backend:", err)
			}
			c.Backends = append(c.Backends, statsd)

		case v == "stdout":
			stdout := NewStdoutBackend()
			c.Backends = append(c.Backends, stdout)

		case v == "librato":
			librato, err := NewLibratoBackend()
			if err != nil {
				log.Fatalln("Impossible to instanciate the Librato backend:", err)
			}
			c.Backends = append(c.Backends, librato)

		case v == "syslog":
			syslog, err := NewSyslogBackend()
			if err != nil {
				log.Fatalln("Impossible to instanciate the syslog backend:", err)
			}
			c.Backends = append(c.Backends, syslog)
		}
	}

	for _, v := range config.Checks {
		check, err := NewCheck(v.Url, v.Key, v.Interval, v.Headers)
		check.Header.Set("User-Agent", config.UserAgent)
		if err != nil {
			log.Fatalln("Check configuration error:", err)
		}
		c.Checks = append(c.Checks, *check)
	}

	c.UserAgent = config.UserAgent
	c.Timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		log.Fatalln("Invalid timeout value given:", err)
	}

	c.Url = os.Getenv("POLLER_URL")
}

func (c *Configuration) CloseBackends() {
	for _, v := range c.Backends {
		v.Close()
	}
}
