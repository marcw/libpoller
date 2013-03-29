package alert

import (
	"github.com/marcw/poller"
)

type Pool map[poller.Alerter]bool

func (p Pool) Add(a poller.Alerter) {
	p[a] = true
}

func (p Pool) Alert(event *poller.Event) {
	for k, _ := range p {
		k.Alert(event)
	}
}
