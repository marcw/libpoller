package poller

// The Config struct holds and links together a CheckList, a Scheduler and a configuration Store.
type Config struct {
	scheduler Scheduler
	store     Store
}

// Instantiates a new Config with an empty CheckList
func NewConfig(store Store, scheduler Scheduler) *Config {
	return &Config{scheduler: scheduler, store: store}
}

func (c *Config) Schedule() {
	c.store.ScheduleAll(c.scheduler)
}

func (c *Config) Add(check *Check) error {
	if err := c.store.Add(check); err != nil {
		return err
	}
	if err := c.scheduler.Schedule(check); err != nil {
		return err
	}

	return nil
}
