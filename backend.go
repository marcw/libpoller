package main

import (
	"fmt"
	"github.com/peterbourgon/g2s"
	"log"
	"os"
	"strconv"
	"time"
)

type Backend interface {
	// "Log" check result in the Backend service
	Log(check *Check, statusCode int, duration time.Duration)
}

// Backend for Statsd
type StatsdBackend struct {
	statsd g2s.Statter
}

// Instanciate a new StatsdBackend
// Uses:
//   STATSD_HOST env variable
//   STATSD_PORT env variable (defaults to 8125)
//   STATSD_PROTOCOL env variable (defaults to udp)
func NewStatsdBackend() (*StatsdBackend, error) {
	envHost := os.Getenv("STATSD_HOST")
	envPort := os.Getenv("STATSD_PORT")
	envProtocol := os.Getenv("STATSD_PROTOCOL")

	if envHost == "" {
		return nil, fmt.Errorf("STATSD_HOST environment variable must be defined")
	}

	if envPort == "" {
		envPort = "8125"
	}

	if envProtocol == "" {
		envProtocol = "udp"
	}

	statsd, err := g2s.Dial(envProtocol, envHost+":"+envPort)
	if err != nil {
		return nil, err
	}

	return &StatsdBackend{statsd}, nil
}

// Log to statsd the check result as follow:
// `Check.Key`.duration : Request duration
// `Check.Key`.success : Request succeeded (status code is 2xx)
// `Check.Key`.error : Request failed (status code != 200)
func (s *StatsdBackend) Log(check *Check, statusCode int, duration time.Duration) {
	s.statsd.Timing(1.0, check.Key+".duration", duration)
	if 200 <= statusCode && 299 >= statusCode {
		s.statsd.Gauge(1.0, check.Key+".success", strconv.Itoa(statusCode))
	} else {
		s.statsd.Gauge(1.0, check.Key+".error", strconv.Itoa(statusCode))
	}
}

// StdoutBackend logs result to Stdout
type StdoutBackend struct {
}

func NewStdoutBackend() *StdoutBackend {

	return &StdoutBackend{}
}

func (s *StdoutBackend) Log(check *Check, statusCode int, duration time.Duration) {
	log.Println(check.Key, statusCode, duration)
}
