package backend

import (
	"fmt"
	"github.com/marcw/poller/check"
	"github.com/rcrowley/go-librato"
	"os"
	"time"
)

type LibratoBackend struct {
	metrics librato.Metrics
	prefix  string
}

func NewLibratoBackend() (*LibratoBackend, error) {
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
		prefix = "poller.checks."
	}

	metrics := librato.NewSimpleMetrics(user, token, source)

	return &LibratoBackend{metrics: metrics, prefix: prefix}, nil
}

func (l *LibratoBackend) Log(e *check.Event) {
	d := l.metrics.GetGauge(l.prefix + e.Check.Key + ".duration")
	d <- int64(e.Duration.Nanoseconds() / int64(time.Millisecond))

	c := l.metrics.GetGauge(l.prefix + e.Check.Key + ".up")
	c <- btou(e.Up)
}

func (l *LibratoBackend) Close() {
	l.metrics.Close()
	l.metrics.Wait()
}
