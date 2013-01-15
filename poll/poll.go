package poll

import (
	"github.com/marcw/poller/check"
)

type Poller interface {
	Poll(c *check.Check) *check.Event
}
