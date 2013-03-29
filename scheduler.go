package poller

import (
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(check *Check)
	Stop(key string)
	StopAll()
	Start()
	Next() <-chan *Check
}

type simpleScheduler struct {
	stopSignals map[string]chan int // collection of channels which are used to signal a goroutine to abandon ship immediately
	toPoll      chan *Check         // checks which are due to polling
	toSchedule  chan *Check         // checks which are due to scheduling
	mu          sync.Mutex
}

// Instantiates a SimpleScheduler which scheduling strategy's fairly basic.
// For each scheduled check, a new time.Timer is created in its own goroutine.
func NewSimpleScheduler() Scheduler {
	return &simpleScheduler{
		stopSignals: make(map[string]chan int),
		toPoll:      make(chan *Check),
		toSchedule:  make(chan *Check)}
}

func (s *simpleScheduler) schedule(check *Check, deleteSignal <-chan int) {
	timer := time.NewTimer(check.Interval)
	select {
	case <-timer.C:
		s.toSchedule <- check
	case <-deleteSignal:
		break
	}
}

func (s *simpleScheduler) Schedule(check *Check) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopSignals[check.Key] = make(chan int)
	go s.schedule(check, s.stopSignals[check.Key])
}

func (s *simpleScheduler) stop(key string) {
	s.stopSignals[key] <- 1
	close(s.stopSignals[key])
	delete(s.stopSignals, key)
}

func (s *simpleScheduler) Stop(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stop(key)
}

func (s *simpleScheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.stopSignals {
		s.stop(k)
	}
}

func (s *simpleScheduler) Start() {
	for {
		check := <-s.toSchedule
		go s.schedule(check, s.stopSignals[check.Key])
		s.toPoll <- check
	}
}

func (s *simpleScheduler) Next() <-chan *Check {
	return s.toPoll
}
