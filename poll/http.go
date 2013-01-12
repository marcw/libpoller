package poll

import (
	"github.com/marcw/poller/check"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpPoller struct {
}

func NewHttpPoller() *HttpPoller {
	return &HttpPoller{}
}

func (p HttpPoller) Poll(c *check.Check) *check.CheckEvent {
	event := check.NewCheckEvent(c)
	timer := time.NewTimer(c.Timeout)

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
