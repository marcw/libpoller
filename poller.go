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

// A Poller is the glue between a Scheduler, a Backend, a Service and a Alerter. 
type Poller interface {
	Run(Scheduler, Backend, Service, Alerter)
}

// A Service is a specialized way to poll a check. ie: a HttpService will specialize in polling HTTP resources.
type Service interface {
	Poll(c *Check) *Event
}

type directPoller struct {
}

// NewDirectPoller() returns a "no-frills" Poller instance. 
// It waits for the next scheduled check, poll it, log it and if alerting is needed, pass it through the alerter.
func NewDirectPoller() Poller {
	return &directPoller{}
}

func (dp *directPoller) Run(scheduler Scheduler, backend Backend, service Service, alerter Alerter) {
	for check := range scheduler.Next() {
		go dp.poll(check, backend, service, alerter)
	}
}

func (db *directPoller) poll(check *Check, backend Backend, service Service, alerter Alerter) {
	event := service.Poll(check)
	go backend.Log(event)
	if event.Check.ShouldAlert() {
		go alerter.Alert(event)
	}
}
