package backend

import (
	"github.com/marcw/poller/check"
	"log"
	"time"
)

// StdoutBackend logs result to Stdout
type StdoutBackend struct {
}

func NewStdoutBackend() *StdoutBackend {

	return &StdoutBackend{}
}

func (s *StdoutBackend) LogSuccess(check *check.Check, statusCode int, duration time.Duration) {
	log.Println(check.Key, statusCode, duration)
}

func (s *StdoutBackend) LogError(check *check.Check, statusCode int, duration time.Duration) {
	log.Println(check.Key, statusCode, duration)
}

func (s *StdoutBackend) LogTimeout(check *check.Check) {
	log.Println(check.Key, "TIMEOUT")
}

func (s *StdoutBackend) Close() {
	// NO OP
}
