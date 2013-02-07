package backend

import (
	"fmt"
	"github.com/marcw/poller/check"
	"github.com/peterbourgon/g2s"
	"os"
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

func (s *StatsdBackend) Log(e *check.Event) {
	s.statsd.Timing(1.0, s.prefix+e.Check.Key+".duration", e.Duration)
	s.statsd.Counter(1.0, s.prefix+e.Check.Key+".up", int(btou(e.IsUp())))
}

func (s *StatsdBackend) Close() {
	// NO OP
}
