package backend

import (
	"github.com/marcw/poller"
	"log"
)

// StdoutBackend logs result to Stdout
type stdoutBackend struct {
}

func NewStdoutBackend() poller.Backend {
	return &stdoutBackend{}
}

func (s *stdoutBackend) Log(e *poller.Event) {
	log.Println(e.Check.Key, btos(e.IsUp()), e.Duration)
}

func (s *stdoutBackend) Close() {
	// NO OP
}
