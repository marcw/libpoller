package check

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func (s *Scheduler) Get(key string) *Check {
	return s.checks[key]
}

func (s *Scheduler) Len() int {
	return len(s.checks)
}

func (s *Scheduler) AddFromJSON(data []byte) error {
	checks := jsonChecks{}
	if err := json.Unmarshal(data, &checks); err != nil {
		return err
	}

	for _, v := range checks {
		chk, err := NewCheck(v.Url, v.Key, v.Interval, v.Headers)
		if err != nil {
			return err
		}

		s.Add(chk)
	}

	return nil
}

func (s *Scheduler) JSON() ([]byte, error) {
	checks := jsonChecks{}
	for _, v := range s.checks {
		header := make(map[string]string)
		for k, h := range v.Header {
			header[k] = h[0]
		}

		check := jsonCheck{Url: v.Url.String(), Key: v.Key, Interval: v.Interval.String(), Headers: header}
		checks = append(checks, check)
	}
	data, err := json.Marshal(checks)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Scheduler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := s.JSON()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Write(data)

		return
	}

	if r.Method == "POST" || r.Method == "PUT" {
		if r.Method == "PUT" {
			s.Wipe()
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)

			return
		}

		err = s.AddFromJSON(data)
		if err != nil {
			http.Error(w, err.Error(), 400)

			return
		}

		w.WriteHeader(201)
		return
	}
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

	for k, _ := range s.checks {
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
