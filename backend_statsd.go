package poller

import (
	"fmt"
	"github.com/peterbourgon/g2s"
	"net"
)

// Backend for Statsd
type statsdBackend struct {
	statsd g2s.Statter
	prefix string
}

// Instanciate a new Backend that will send data to a statsd instance
func NewStatsdBackend(host, port, protocol, prefix string) (Backend, error) {
	if host == "" {
		return nil, fmt.Errorf("Statsd host cannot be empty")
	}

	if port == "" {
		port = "8125"
	}

	if protocol == "" {
		protocol = "udp"
	}

	if prefix == "" {
		prefix = "checks."
	}

	statsd, err := g2s.Dial(protocol, net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}

	return &statsdBackend{statsd: statsd, prefix: prefix}, nil
}

func (s *statsdBackend) Log(e *Event) {
	s.statsd.Timing(1.0, s.prefix+e.Check.Key+".duration", e.Duration)
	s.statsd.Counter(1.0, s.prefix+e.Check.Key+".up", int(btou(e.IsUp())))
}

func (s *statsdBackend) Close() {
	// NO OP
}
