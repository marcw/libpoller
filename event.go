package poller

import (
	"time"
)

// Represents the state of a Check after being polled
type Event struct {
	Check      *Check        // check
	Duration   time.Duration // total duration of check
	StatusCode int           // http status code, if any
	Time       time.Time     // time of check
	up         bool          // true if service is up
	Alert      bool          // true if backend should raise an alert
	NotifyFix  bool          // true if backend should notify of service being up again
}

func NewEvent(check *Check) *Event {
	return &Event{Time: time.Now(), Check: check}
}

func (e *Event) IsUp() bool {
	return e.up
}

func (e *Event) Up() {
	e.up = true

	if e.Check.UpSince.IsZero() {
		e.Check.UpSince = e.Time
	}

	// if service was down
	if !e.Check.DownSince.IsZero() {
		e.Check.UpSince = e.Time
		e.Check.WasDownFor = e.Time.Sub(e.Check.DownSince)
		e.Check.DownSince = time.Time{}
		e.NotifyFix = e.Check.NotifyFix
	}
}

func (e *Event) Down() {
	e.up = false

	if e.Check.DownSince.IsZero() {
		e.Check.DownSince = e.Time
	}

	// if service was up
	if !e.Check.UpSince.IsZero() {
		e.Check.DownSince = e.Time
		e.Check.WasUpFor = e.Time.Sub(e.Check.UpSince)
		e.Check.UpSince = time.Time{}

		// Is it time we alert backend?
		if !e.Check.ShouldAlert() {
			e.Alert = true
			e.Check.Alerted = true
		}
	}
}
