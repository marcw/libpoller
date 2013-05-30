package poller

import (
	"fmt"
	"github.com/rcrowley/go-librato"
	"time"
)

type libratoBackend struct {
	metrics librato.Metrics
	prefix  string
}

func NewLibratoBackend(user, token, source, prefix string) (Backend, error) {
	if user == "" {
		return nil, fmt.Errorf("Librato user cannot be empty")
	}

	if token == "" {
		return nil, fmt.Errorf("Librato token cannot be empty")
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
