package poll

import (
	"github.com/marcw/poller/check"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpPoller struct {
    UserAgent string
    Timeout   time.Duration
}

func NewHttpPoller(ua string, timeout time.Duration) *HttpPoller {
	return &HttpPoller{UserAgent: ua, Timeout: timeout}
}

func (p HttpPoller) Poll(c *check.Check) *check.CheckEvent {
	event := check.NewCheckEvent(c)
	timer := time.NewTimer(p.Timeout)

	start := time.Now().UnixNano()
	conn, err := net.Dial("tcp", c.Addr.String())
	if err != nil {
		end := time.Now().UnixNano()
		event.Duration = time.Duration(end - start)

		return event
	}

	defer conn.Close()
	ch := make(chan *check.CheckEvent, 1)
	go func() {
		client := httputil.NewClientConn(conn, nil)
		req, err := http.NewRequest("GET", c.Url.String(), nil)
		req.Header = c.Header
        req.Header.Set("User-Agent", p.UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			ch <- event
			return
		}
		defer resp.Body.Close()

		end := time.Now().UnixNano()

		event.StatusCode = resp.StatusCode
		event.Duration = time.Duration(end - start)
		if event.StatusCode == 200 {
			event.Up = true
		}

		ch <- event
	}()

	select {
	case <-timer.C:
		end := time.Now().UnixNano()
		event.Duration = time.Duration(end - start)
		event.Timeout = true
		event.Up = false

		return event

	case e := <-ch:
		return e
	}
	panic("unreachable")
}
