package poller

// An Alerter raises an alert based on the event it received.
// An alert is a communication to a system or a user with the information about current's and past check's states.
// For concrete implementation, see the "github.com/marcw/poller/alert" package.
type Alerter interface {
	Alert(event *Event)
}

// A backend log checks event.
// For concrete implementation, see the "github.com/marcw/poller/backend" package.
type Backend interface {
	Log(e *Event)
	Close()
}

// A Poller is the glue between a Scheduler, a Backend, a Probe and a Alerter.
type Poller interface {
	Run(Scheduler, Backend, Probe, Alerter)
}

// A Probe is a specialized way to poll a check. ie: a HttpProbe will specialize in polling HTTP resources.
type Probe interface {
	Test(c *Check) *Event
}

// A Store defines a place where configuration can be loaded/persisted.
type Store interface {
	Add(*Check) error
	Get(key string) (*Check, error)
	Remove(key string) error
	Len() (int, error)
	ScheduleAll(Scheduler) error
}

type directPoller struct {
}

// NewDirectPoller() returns a "no-frills" Poller instance.
// It waits for the next scheduled check, poll it, log it and if alerting is needed, pass it through the alerter.
func NewDirectPoller() Poller {
	return &directPoller{}
}

func (dp *directPoller) Run(scheduler Scheduler, backend Backend, probe Probe, alerter Alerter) {
	for check := range scheduler.Next() {
		go dp.poll(check, backend, probe, alerter)
	}
}

func (db *directPoller) poll(check *Check, backend Backend, probe Probe, alerter Alerter) {
	event := probe.Test(check)
	go backend.Log(event)
	if event.Check.ShouldAlert() {
		go alerter.Alert(event)
	}
}
