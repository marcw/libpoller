package poll

import (
	"github.com/marcw/poller/check"
	"time"
)

type Poller interface {
	Poll(c *check.Check) (int, time.Duration, error)
}
