package main

import (
	"encoding/json"
	"log"
)

type Configuration struct {
	Backends []Backend
	Checks   []Check
}

func (c *Configuration) Load(data []byte) {
	type check struct {
		Url      string
		Key      string
		Interval string
	}

	type configuration struct {
		Backends []string
		Checks   []check
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
		}
	}

	for _, v := range config.Checks {
		check, err := NewCheck(v.Url, v.Key, v.Interval)
		if err != nil {
			log.Fatalln("Check configuration error:", err)
		}
		c.Checks = append(c.Checks, *check)
	}
}
