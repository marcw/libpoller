package poller

type Alerter interface {
	Alert(event *Event)
}

type Backend interface {
	Log(e *Event)
	Close()
}

type Poller interface {
	Run(Scheduler, Backend, Service, Alerter)
}

type Service interface {
	Poll(c *Check) *Event
}

type directPoller struct {
}

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
