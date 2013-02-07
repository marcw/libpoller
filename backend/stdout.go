package backend

import (
	"github.com/marcw/poller/check"
	"log"
)

// StdoutBackend logs result to Stdout
type StdoutBackend struct {
}

func NewStdoutBackend() *StdoutBackend {
	return &StdoutBackend{}
}

func (s *StdoutBackend) Log(e *check.Event) {
	log.Println(e.Check.Key, btos(e.IsUp()), e.Duration)
	if e.Alert {
		log.Println(e.Check.Key, "ALERT", "Down since", e.Check.DownSince)
	}
}

func (s *StdoutBackend) Close() {
	// NO OP
}
