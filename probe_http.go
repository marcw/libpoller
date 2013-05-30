package poller

import (
	"net/http"
	"time"
)

type httpProbe struct {
	UserAgent string
	Timeout   time.Duration
}

func NewHttpProbe(ua string, timeout time.Duration) Probe {
	return &httpProbe{UserAgent: ua, Timeout: timeout}
}

func (p *httpProbe) Test(c *Check) *Event {
	event := NewEvent(c)
	timer := time.NewTimer(p.Timeout)
	ch := make(chan *Event, 1)

	start := time.Now().UnixNano()
	go func(e *Event, eventCh chan<- *Event) {
		client := &http.Client{Jar: nil}
		req, err := http.NewRequest("GET", c.Config.GetString("url"), nil)
		var header = http.Header{}

		for k, v := range c.Config.GetMapStringString("headers") {
			header.Set(k, v)
		}

		req.Header = header
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
}
