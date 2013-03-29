package check

import (
	"sync"
	"time"
)

// Contains a collection of Check
type Scheduler struct {
	checks        map[string]*Check   // collection of checks
	deleteSignals map[string]chan int // collection of channels which are used to signal a goroutine to abandon ship immediately
	ToPoll        chan *Check         // checks which are due to polling
	toSchedule    chan *Check         // checks which are due to scheduling
	w             sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		checks:        make(map[string]*Check),
		deleteSignals: make(map[string]chan int),
		ToPoll:        make(chan *Check),
		toSchedule:    make(chan *Check)}
}

func (s *Scheduler) runCheck(c *Check, deleteSignal <-chan int) {
	timer := time.NewTimer(c.Interval)
	select {
	case <-timer.C:
		s.toSchedule <- c
	case <-deleteSignal:
		break
	}
}

func (s *Scheduler) Add(c *Check) {
	s.w.Lock()
	defer s.w.Unlock()

	s.checks[c.Key] = c
	s.deleteSignals[c.Key] = make(chan int)
	go s.runCheck(c, s.deleteSignals[c.Key])
}

func (s *Scheduler) Delete(key string) {
	s.w.Lock()
	defer s.w.Unlock()

	s.del(key)
}

func (s *Scheduler) del(key string) {
	s.deleteSignals[key] <- 0
	delete(s.checks, key)
	delete(s.deleteSignals, key)
}

func (s *Scheduler) Wipe() {
	s.w.Lock()
	defer s.w.Unlock()

	for k := range s.checks {
		s.del(k)
	}
}

func (s *Scheduler) Run() {
	for {
		check := <-s.toSchedule
		go s.runCheck(check, s.deleteSignals[check.Key])
		s.ToPoll <- check
	}
}
