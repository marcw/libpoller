package poller

import (
	"encoding/json"
	"sync"
)

// A CheckList contains a list of checks. The underlying data structure is a map[string]*Check
// Concurrent access in read and write is protected by a mutex.
type CheckList struct {
	list map[string]*Check
	mu   sync.Mutex
}

// Instantiates a new CheckList
func NewCheckList() *CheckList {
	return &CheckList{list: make(map[string]*Check)}
}

// Instantiates a new CheckList and populates it with JSON
func NewCheckListFromJSON(data []byte) (*CheckList, error) {
	cl := NewCheckList()
	var checks []*jsonCheck
	if err := json.Unmarshal(data, &checks); err != nil {
		return nil, err
	}

	for _, v := range checks {
		chk, err := v.toCheck()
		if err != nil {
			return nil, err
		}

		cl.Add(chk)
	}

	return cl, nil
}

// Add an element to the list
func (c *CheckList) Add(check *Check) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list[check.Key] = check
}

// Returns the element key from the list. If key is not present, nil is returned.
func (c *CheckList) Get(key string) *Check {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.list[key]
}

// Delete removes element identified by key from the list.
func (c *CheckList) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.list, key)
}

// Len returns the number of items in the CheckList
func (c *CheckList) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.list)
}

// Each will apply function `each` to each checks of the list.
// During the whole time of the execution, the CheckList is Locked
// for both reading and writing.
func (c *CheckList) Each(each func(check *Check)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.list {
		each(c.list[k])
	}
}

// JSON() marshals the content of the CheckList as JSON.
func (c *CheckList) JSON() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var checks []*jsonCheck
	for k := range c.list {
		checks = append(checks, c.list[k].json())
	}

	data, err := json.Marshal(checks)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Clear all elements of the list
func (c *CheckList) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list = make(map[string]*Check)
}
