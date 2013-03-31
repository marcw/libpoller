package alert

import (
	"fmt"
	"github.com/marcw/pagerduty"
	"github.com/marcw/poller"
	"os"
	"time"
)

type pagerDutyAlerter struct {
	serviceKey string
}

func NewPagerDutyAlerter() (poller.Alerter, error) {
	envServiceKey := os.Getenv("PAGERDUTY_SERVICE_KEY")
	if envServiceKey == "" {
		return nil, fmt.Errorf("Please define the PAGERDUTY_SERVICE_KEY environment variable.")
	}

	return &pagerDutyAlerter{envServiceKey}, nil
}

func (pda *pagerDutyAlerter) Alert(event *poller.Event) {
	description := fmt.Sprintf("%s (%s) is DOWN since %s.", event.Check.Key, event.Check.Url.String(), event.Check.DownSince.Format(time.RFC3339))
	e := pagerduty.NewTriggerEvent(pda.serviceKey, description)
	e.Details["checked_at"] = event.Time.Format(time.RFC3339)
	e.Details["duration"] = event.Duration.String()
	e.Details["status_code"] = event.StatusCode
	e.Details["was_up_for"] = event.Check.WasUpFor.String()
	e.IncidentKey = event.Check.Key
	for {
		_, statusCode, _ := pagerduty.Submit(e)
		if statusCode < 500 {
			break
		} else {
			// Wait a bit before trying again
			time.Sleep(3 * time.Second)
		}
	}
}
