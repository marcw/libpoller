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
