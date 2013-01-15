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
	log.Println(e.Check.Key, btos(e.Up), e.Duration)
}

func (s *StdoutBackend) Close() {
	// NO OP
}
