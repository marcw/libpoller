package poller

// The Config struct holds and links together a CheckList, a Scheduler and a configuration Store.
type Config struct {
	checks    *CheckList
	scheduler Scheduler
	store     Store
}

// Instantiates a new Config with an empty CheckList
func NewConfig(store Store, scheduler Scheduler) *Config {
	return &Config{checks: NewCheckList(), scheduler: scheduler, store: store}
}

// Reload the entire configuration from the store, wipe the scheduler's content
// and replace it with the new checks.
func (c *Config) Load() error {
	if checks, err := c.store.Load(); err != nil {
		return err
	} else {
		c.checks = checks
	}

	c.scheduler.StopAll()
	c.checks.Each(func(check *Check) {
		c.scheduler.Schedule(check)
	})

	return nil
}

// Persist the configuration changes since last update. Only changed checks from last updates will be passed to the store.
func (c *Config) Persist() {
	c.store.Persist(c.checks)
}

// Add a new check to the CheckList and schedule it
func (c *Config) Add(chk *Check) {
	c.checks.Add(chk)
	c.scheduler.Schedule(chk)
}

func (c *Config) Clear() {
	c.checks.Clear()
	c.scheduler.StopAll()
}

func (c *Config) SetCheckList(cl *CheckList) {
	c.Clear()
	c.checks = cl
	cl.Each(func(check *Check) {
		c.scheduler.Schedule(check)
	})
}

func (c *Config) Scheduler() Scheduler {
	return c.scheduler
}
