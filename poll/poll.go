package poll

import (
	"github.com/marcw/poller"
)

type Poller interface {
	Poll(c *poller.Check) *poller.Event
}
