package main

import (
	"fmt"
	"github.com/peterbourgon/g2s"
	"github.com/rcrowley/go-librato"
	"log"
	"os"
	"time"
)

type Backend interface {
	// "Log" check result in the Backend service
	Log(check *Check, statusCode int, duration time.Duration)
	Close()
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
		s.statsd.Counter(1.0, check.Key+".success", 1)
	} else {
		s.statsd.Counter(1.0, check.Key+".error", 1)
	}
}

func (s *StatsdBackend) Close() {
	// NO OP
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

func (s *StdoutBackend) Close() {
	// NO OP
}

type LibratoBackend struct {
	metrics librato.Metrics
}

func NewLibratoBackend() (*LibratoBackend, error) {
	user := os.Getenv("LIBRATO_USER")
	token := os.Getenv("LIBRATO_TOKEN")
	source := os.Getenv("LIBRATO_SOURCE")

	if user == "" {
		return nil, fmt.Errorf("LIBRATO_USER environment variable must be defined")
	}

	if token == "" {
		return nil, fmt.Errorf("LIBRATO_TOKEN environment variable must be defined")
	}

	if source == "" {
		source = "poller"
	}

	metrics := librato.NewSimpleMetrics(user, token, source)

	return &LibratoBackend{metrics}, nil
}

func (l *LibratoBackend) Log(check *Check, statusCode int, duration time.Duration) {
	d := l.metrics.GetGauge(check.Key + ".duration")
	d <- int64(duration.Nanoseconds() / int64(time.Millisecond))
	if 200 <= statusCode && 299 >= statusCode {
		c := l.metrics.GetCounter(check.Key + ".success")
		c <- 1
	} else {
		c := l.metrics.GetCounter(check.Key + ".error")
		c <- 1
	}
}

func (l *LibratoBackend) Close() {
	l.metrics.Close()
	l.metrics.Wait()
}
