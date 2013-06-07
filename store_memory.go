package poller

import (
	"sync"
)

// A CheckList contains a list of checks. The underlying data structure is a map[string]*Check
// Concurrent access in read and write is protected by a mutex.
type inMemoryStore struct {
	list map[string]*Check
	mu   sync.Mutex
}

// Instantiates a new CheckList
func NewInMemoryStore() Store {
	return &inMemoryStore{list: make(map[string]*Check)}
}

// Add an element to the list
func (s *inMemoryStore) Add(check *Check) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list[check.Key] = check
	return nil
}

// Returns the element key from the list. If key is not present, nil is returned.
func (s *inMemoryStore) Get(key string) (*Check, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.list[key], nil
}

// Delete removes element identified by key from the list.
func (s *inMemoryStore) Remove(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.list, key)
	return nil
}

// Len returns the number of items in the CheckList
func (s *inMemoryStore) Len() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.list), nil
}

func (s *inMemoryStore) ScheduleAll(scheduler Scheduler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.list {
		scheduler.Schedule(s.list[k])
	}

	return nil
}
