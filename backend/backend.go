package backend

import (
	"github.com/marcw/poller"
)

type Backend interface {
	Log(e *poller.Event)
	Close()
}

type Pool map[Backend]bool

func (p Pool) Add(b Backend) {
	p[b] = true
}

func (p Pool) Log(event *poller.Event) {
	for k, _ := range p {
		k.Log(event)
	}
}

func (p Pool) Close() {
	for k, _ := range p {
		k.Close()
	}
}

func btou(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func btos(b bool) string {
	if b {
		return "UP"
	}
	return "DOWN"
}
