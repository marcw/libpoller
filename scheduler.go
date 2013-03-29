package poller

import (
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(check *Check)
	Stop(key string)
	StopAll()
	Next() <-chan *Check
}

// Contains a collection of Check
type SimpleScheduler struct {
	stopSignals map[string]chan int // collection of channels which are used to signal a goroutine to abandon ship immediately
	toPoll      chan *Check         // checks which are due to polling
	toSchedule  chan *Check         // checks which are due to scheduling
	mu          sync.Mutex
}

// Instantiates a SimpleScheduler which scheduling strategy's fairly basic.
// For each scheduled check, a new time.Timer is created in its own goroutine.
func NewSimpleScheduler() *SimpleScheduler {
	return &SimpleScheduler{
		stopSignals: make(map[string]chan int),
		toPoll:      make(chan *Check),
		toSchedule:  make(chan *Check)}
}

func (s *SimpleScheduler) schedule(check *Check, deleteSignal <-chan int) {
	timer := time.NewTimer(check.Interval)
	select {
	case <-timer.C:
		s.toSchedule <- check
	case <-deleteSignal:
		break
	}
}

func (s *SimpleScheduler) Schedule(check *Check) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopSignals[check.Key] = make(chan int)
	go s.schedule(check, s.stopSignals[check.Key])
}

func (s *SimpleScheduler) stop(key string) {
	s.stopSignals[key] <- 1
	close(s.stopSignals[key])
	delete(s.stopSignals, key)
}

func (s *SimpleScheduler) Stop(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stop(key)
}

func (s *SimpleScheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.stopSignals {
		s.stop(k)
	}
}

func (s *SimpleScheduler) Run() {
	for {
		check := <-s.toSchedule
		go s.schedule(check, s.stopSignals[check.Key])
		s.toPoll <- check
	}
}

func (s *SimpleScheduler) Next() <-chan *Check {
	return s.toPoll
}
