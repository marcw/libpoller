package alert

import (
	"github.com/marcw/poller"
)

type Alerter interface {
	Alert(event *poller.Event)
}

type Pool map[Alerter]bool

func (p Pool) Add(a Alerter) {
	p[a] = true
}

func (p Pool) Alert(event *poller.Event) {
	for k, _ := range p {
		k.Alert(event)
	}
}
