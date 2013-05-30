package poller

import (
	"fmt"
	"github.com/marcw/bag"
	"time"
)

type CheckType string

const (
	CheckTypeUDP  CheckType = "udp"
	CheckTypeHTTP CheckType = "http"
)

type Check struct {
	Key       string    // Key (should be unique among same Scheduler
	checkType CheckType // Type of check

	Interval time.Duration // Interval between each check

	UpSince    time.Time     // Time since the service is up
	DownSince  time.Time     // Time since the service is down
	WasDownFor time.Duration // Time since the service was down
	WasUpFor   time.Duration // Time since the service was up

	Alert     bool // Raise alert if service is down
	Alerted   bool // Is backend already alerted?
	NotifyFix bool // Notify if service is back up

	AlertDelay time.Duration // Delay before raising an alert (zero value = NOW)
	Config     *bag.Bag
}

func newCheck() *Check {
	return &Check{Config: bag.NewBag()}
}

func NewCheck(key, interval string, alert bool, alertDelay string, notifyFix bool, config map[string]interface{}) (*Check, error) {
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

	return &Check{Key: key, Interval: d, Alert: alert, AlertDelay: ad, NotifyFix: notifyFix, Config: bag.From(config)}, nil
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

func (c *Check) Type() CheckType {
	return c.checkType
}

func (c *Check) AlertDescription() string {
	if c.Type() == CheckTypeHTTP {
		return fmt.Sprintf("%s (%s)", c.Key, c.Config.GetString("url"))
	}

	return ""
}
