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

func (p HttpPoller) Poll(c *check.Check) (int, time.Duration, error) {
	conn, err := net.DialTimeout("tcp", c.Url.Host, 10*time.Second)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()
	client := httputil.NewClientConn(conn, nil)
	req, err := http.NewRequest("GET", c.Url.String(), nil)
	req.Header = c.Header
	start := time.Now().UnixNano()
	resp, err := client.Do(req)
	end := time.Now().UnixNano()
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, time.Duration(end - start), nil
}
