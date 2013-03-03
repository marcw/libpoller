package alert

import (
	"github.com/marcw/poller/check"
)

type Alerter interface {
	Alert(event *check.Event)
}

type Pool map[Alerter]bool

func (p Pool) Add(a Alerter) {
	p[a] = true
}

func (p Pool) Alert(event *check.Event) {
	for k, _ := range p {
		k.Alert(event)
	}
}
