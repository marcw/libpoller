package poller

import (
	"net"
	"strconv"
	"time"
)

type udpProbe struct {
	Timeout time.Duration
}

func NewUdpProbe(timeout time.Duration) Probe {
	return &udpProbe{timeout}
}

func (p *udpProbe) Test(c *Check) *Event {
	event := NewEvent(c)
	timer := time.NewTimer(p.Timeout)

	ch := make(chan *Event, 1)

	start := time.Now().UnixNano()
	go func(e *Event, eventCh chan<- *Event) {
		hp := net.JoinHostPort(c.Config.GetString("host"), strconv.Itoa(c.Config.GetInt("port")))
		raddr, err := net.ResolveUDPAddr("udp", hp)
		if err != nil {
			eventCh <- e
			return
		}
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			eventCh <- e
			return
		}
		conn.SetDeadline(time.Now().Add(p.Timeout))
		defer conn.Close()

		if _, err := conn.Write([]byte(c.Config.GetString("send"))); err != nil {
			eventCh <- e
			return
		}
		buf := make([]byte, len([]byte(c.Config.GetString("receive"))))
		for {
			count, err := conn.Read(buf)
			if err != nil {
				eventCh <- e
				return
			}
			if count != 0 {
				break
			}
		}
		if string(buf) == c.Config.GetString("receive") {
			e.Up()
			eventCh <- e
			return
		} else {
			e.Down()
			eventCh <- e
			return
		}
	}(event, ch)

	select {
	case <-timer.C:
		end := time.Now().UnixNano()
		event.Duration = time.Duration(end - start)
		event.Down()

		return event

	case e := <-ch:
		end := time.Now().UnixNano()
		event.Duration = time.Duration(end - start)
		return e
	}
}
