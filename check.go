package poller

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type Check struct {
	Url        *url.URL      // URL of check
	Addr       net.Addr      // 
	Key        string        // Key (should be unique among same Scheduler
	Interval   time.Duration // Interval between each check
	Header     http.Header   // HTTP Headers (if any)
	UpSince    time.Time     // Time since the service is up
	DownSince  time.Time     // Time since the service is down
	WasDownFor time.Duration // Time since the service was down
	WasUpFor   time.Duration // Time since the service was up
	Alert      bool          // Raise alert if service is down
	Alerted    bool          // Is backend already alerted?
	NotifyFix  bool          // Notify if service is back up
	AlertDelay time.Duration // Delay before raising an alert (zero value = NOW)
}

func NewCheck(checkUrl, key, interval string, alert bool, alertDelay string, notifyFix bool, headers map[string]string) (*Check, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	var ad time.Duration
	if alert {
		ad, err = time.ParseDuration(alertDelay)
		if err != nil {
			return nil, err
		}
	}

	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}

	u, err := url.Parse(checkUrl)
	if err != nil {
		return nil, err
	}

	host := u.Host
	_, err = net.ResolveTCPAddr("tcp", u.Host)
	if err != nil {
		if u.Scheme == "http" {
			host = host + ":80"
		} else {
			host = host + ":443"
		}
	}

	a, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}

	return &Check{Url: u, Key: key, Interval: d, Header: h, Addr: a, Alert: alert, AlertDelay: ad, NotifyFix: notifyFix}, nil
}

// Check if it's time to send the alert. Returns true if it is.
func (c *Check) ShouldAlert() bool {
	return c.Alert && !c.Alerted && c.DownSince.Add(c.AlertDelay).Before(time.Now())
}

func (c *Check) ShouldNotifyFix() bool {
	if !c.NotifyFix {
		return false
	}

	if c.WasDownFor == 0 {
		return false
	}

	if c.WasDownFor > 0 && c.NotifyFix && c.Alerted {
		return true
	}

	return false
}
