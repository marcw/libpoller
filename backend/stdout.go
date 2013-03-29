package backend

import (
	"github.com/marcw/poller"
	"log"
)

// StdoutBackend logs result to Stdout
type StdoutBackend struct {
}

func NewStdoutBackend() *StdoutBackend {
	return &StdoutBackend{}
}

func (s *StdoutBackend) Log(e *poller.Event) {
	log.Println(e.Check.Key, btos(e.IsUp()), e.Duration)
}

func (s *StdoutBackend) Close() {
	// NO OP
}
