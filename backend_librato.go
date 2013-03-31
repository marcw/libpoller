package poller

import (
	"fmt"
	"github.com/rcrowley/go-librato"
	"os"
	"time"
)

type libratoBackend struct {
	metrics librato.Metrics
	prefix  string
}

func NewLibratoBackend() (Backend, error) {
	user := os.Getenv("LIBRATO_USER")
	token := os.Getenv("LIBRATO_TOKEN")
	source := os.Getenv("LIBRATO_SOURCE")
	prefix := os.Getenv("LIBRATO_PREFIX")

	if user == "" {
		return nil, fmt.Errorf("LIBRATO_USER environment variable must be defined")
	}

	if token == "" {
		return nil, fmt.Errorf("LIBRATO_TOKEN environment variable must be defined")
	}

	if source == "" {
		source = "poller"
	}

	if prefix == "" {
		prefix = "checks."
	}

	metrics := librato.NewSimpleMetrics(user, token, source)

	return &libratoBackend{metrics: metrics, prefix: prefix}, nil
}

func (l *libratoBackend) Log(e *Event) {
	d := l.metrics.GetGauge(l.prefix + e.Check.Key + ".duration")
	d <- int64(e.Duration.Nanoseconds() / int64(time.Millisecond))

	c := l.metrics.GetGauge(l.prefix + e.Check.Key + ".up")
	c <- btou(e.IsUp())
}

func (l *libratoBackend) Close() {
	l.metrics.Close()
	l.metrics.Wait()
}
