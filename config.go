package poller

import (
	"encoding/json"
	"github.com/marcw/poller/check"
	"io/ioutil"
	"net/http"
)

type CheckList map[string]*check.Check

// A Store defines a place where configuration can be loaded/persisted.
type Store interface {
	// Read configuration and returns the complete list of checks.
	Load() (CheckList, error)

	// Persist `checks`.
	Persist(CheckList) error
}

type InMemoryStore struct {
	checks CheckList
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{make(CheckList)}
}

func (ims *InMemoryStore) Load() (CheckList, error) {
	return ims.checks, nil
}

func (ims *InMemoryStore) Persist(checks CheckList) error {
	ims.checks = checks
	return nil
}

type Config struct {
	checks    map[string]*check.Check
	scheduler *check.Scheduler
	store     Store
}

// Used for marshalling / unmarshalling
type jsonCheck struct {
	Url        string
	Key        string
	Interval   string
	Alert      bool
	AlertDelay string
	NotifyFix  bool
	Headers    map[string]string
}

type jsonChecks []jsonCheck

func NewConfig(store Store, scheduler *check.Scheduler) *Config {
	return &Config{checks: make(CheckList), scheduler: scheduler, store: store}
}

// Reload the entire configuration from the store, wipe the scheduler's content
// and replace it with the new checks.
func (c *Config) Load() error {
	if checks, err := c.store.Load(); err != nil {
		return err
	} else {
		c.checks = checks
	}

	c.scheduler.Wipe()
	for k := range c.checks {
		c.scheduler.Add(c.checks[k])
	}

	return nil
}

// Persist the configuration changes since last update. Only changed checks from last updates will be passed to the store.
func (c *Config) Persist() {
	c.store.Persist(c.checks)
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := c.JSON()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		w.Write(data)

		return
	}

	if r.Method == "POST" || r.Method == "PUT" {
		if r.Method == "PUT" {
			c.scheduler.Wipe()
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)

			return
		}

		err = c.AddFromJSON(data)
		if err != nil {
			http.Error(w, err.Error(), 400)

			return
		}

		w.WriteHeader(201)
		return
	}
}

func (c *Config) JSON() ([]byte, error) {
	checks := jsonChecks{}
	for _, v := range c.checks {
		header := make(map[string]string)
		for k, h := range v.Header {
			header[k] = h[0]
		}

		check := jsonCheck{Url: v.Url.String(), Key: v.Key, Interval: v.Interval.String(), Alert: v.Alert, AlertDelay: v.AlertDelay.String(), Headers: header}
		checks = append(checks, check)
	}
	data, err := json.Marshal(checks)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Config) AddFromJSON(data []byte) error {
	checks := jsonChecks{}
	if err := json.Unmarshal(data, &checks); err != nil {
		return err
	}

	for _, v := range checks {
		chk, err := check.NewCheck(v.Url, v.Key, v.Interval, v.Alert, v.AlertDelay, v.NotifyFix, v.Headers)
		if err != nil {
			return err
		}

		c.Add(chk)
	}

	return nil
}

func (c *Config) Add(chk *check.Check) {
	c.checks[chk.Key] = chk
	c.scheduler.Add(chk)
}

func (c *Config) Get(key string) *check.Check {
	return c.checks[key]
}

func (c *Config) Len() int {
	return len(c.checks)
}

func (c *Config) Clear() {
	c.checks = make(CheckList)
}
