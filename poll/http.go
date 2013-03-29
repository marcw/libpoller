package poll

import (
	"github.com/marcw/poller/check"
	"net/http"
	"time"
)

type HttpPoller struct {
	UserAgent string
	Timeout   time.Duration
}

func NewHttpPoller(ua string, timeout time.Duration) *HttpPoller {
	return &HttpPoller{UserAgent: ua, Timeout: timeout}
}

func (p *HttpPoller) Poll(c *check.Check) *check.Event {
	event := check.NewEvent(c)
	timer := time.NewTimer(p.Timeout)
	ch := make(chan *check.Event, 1)

	start := time.Now().UnixNano()
	go func(e *check.Event, eventCh chan<- *check.Event) {
		client := &http.Client{Jar: nil}
		req, err := http.NewRequest("GET", c.Url.String(), nil)
		req.Header = c.Header
		req.Header.Set("User-Agent", p.UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			eventCh <- e
			return
		}
		defer resp.Body.Close()

		e.StatusCode = resp.StatusCode
		if e.StatusCode == 200 {
			e.Up()
		} else {
			e.Down()
		}

		eventCh <- e
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
	panic("unreachable")
}
