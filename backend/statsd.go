package backend

import (
	"fmt"
	"github.com/marcw/poller/check"
	"github.com/peterbourgon/g2s"
	"os"
	"time"
)

// Backend for Statsd
type StatsdBackend struct {
	statsd g2s.Statter
	prefix string
}

// Instanciate a new StatsdBackend
// Uses:
//   STATSD_HOST env variable
//   STATSD_PORT env variable (defaults to 8125)
//   STATSD_PROTOCOL env variable (defaults to udp)
//   STATSD_PREFIX env variable (defaults to `poller.checks.`)
func NewStatsdBackend() (*StatsdBackend, error) {
	envHost := os.Getenv("STATSD_HOST")
	envPort := os.Getenv("STATSD_PORT")
	envProtocol := os.Getenv("STATSD_PROTOCOL")
	envPrefix := os.Getenv("STATSD_PREFIX")

	if envHost == "" {
		return nil, fmt.Errorf("STATSD_HOST environment variable must be defined")
	}

	if envPort == "" {
		envPort = "8125"
	}

	if envProtocol == "" {
		envProtocol = "udp"
	}

	if envPrefix == "" {
		envPrefix = "poller.checks."
	}

	statsd, err := g2s.Dial(envProtocol, envHost+":"+envPort)
	if err != nil {
		return nil, err
	}

	return &StatsdBackend{statsd: statsd, prefix: envPrefix}, nil
}

// Log to statsd the check result as follow:
// `check.Check.Key`.duration : Request duration
// `check.Check.Key`.success : Request succeeded (status code is 2xx)
// `check.Check.Key`.error : Request failed (status code != 200)
func (s *StatsdBackend) LogSuccess(check *check.Check, statusCode int, duration time.Duration) {
	s.logDuration(check, duration)
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 1)
}

func (s *StatsdBackend) LogError(check *check.Check, statusCode int, duration time.Duration) {
	s.logDuration(check, duration)
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 0)
}

func (s *StatsdBackend) LogTimeout(check *check.Check) {
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 0)
}

func (s *StatsdBackend) logDuration(check *check.Check, duration time.Duration) {
	s.statsd.Timing(1.0, s.prefix+check.Key+".duration", duration)
}

func (s *StatsdBackend) Close() {
	// NO OP
}
