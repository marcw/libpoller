package poller

import (
	"log"
)

// StdoutBackend logs result to Stdout
type stdoutBackend struct {
}

func NewStdoutBackend() Backend {
	return &stdoutBackend{}
}

func (s *stdoutBackend) Log(e *Event) {
	log.Println(e.Check.Key, btos(e.IsUp()), e.Duration)
}

func (s *stdoutBackend) Close() {
	// NO OP
}
