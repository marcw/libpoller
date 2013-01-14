package backend

import (
	"github.com/marcw/poller/check"
	"time"
)

type Backend interface {
	LogSuccess(check *check.Check, statusCode int, duration time.Duration)
	LogError(Check *check.Check, statusCode int, duration time.Duration)
	LogTimeout(check *check.Check)
	Close()
}

type BackendsList map[Backend]bool

var Backends = make(BackendsList)

func (bl BackendsList) Add(b Backend) {
    bl[b] = true
}

func (bl BackendsList) LogSuccess(c *check.Check, statusCode int, duration time.Duration) {
	for k, _ := range bl {
        k.LogSuccess(c, statusCode, duration)
	}
}

func (bl BackendsList) LogError(c *check.Check, statusCode int, duration time.Duration) {
	for k, _ := range bl {
        k.LogError(c, statusCode, duration)
	}
}

func (bl BackendsList) LogTimeout(c *check.Check) {
	for k, _ := range bl {
        k.LogTimeout(c)
	}
}

func (bl BackendsList) Close() {
	for k, _ := range bl {
        k.Close()
	}
}

// Will init backends based on their defaults environment value
func Add(b Backend) {
    Backends.Add(b)
}

