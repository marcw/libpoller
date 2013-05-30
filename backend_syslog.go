package poller

import (
	"fmt"
	"log/syslog"
)

type syslogBackend struct {
	writer *syslog.Writer
}

// Create a new SyslogBackend instance.
func NewSyslogBackend(network, raddr, prefix string) (Backend, error) {
	if prefix == "" {
		prefix = "poller"
	}

	writer, err := syslog.Dial(network, raddr, syslog.LOG_INFO, prefix)
	if err != nil {
		return nil, err
	}

	return &syslogBackend{writer: writer}, nil
}

func (s *syslogBackend) Log(e *Event) {
	if e.IsUp() {
		s.writer.Info(fmt.Sprintln(e.Check.Key, btos(e.IsUp()), e.Duration))
	} else {
		s.writer.Err(fmt.Sprintln(e.Check.Key, btos(e.IsUp()), e.Duration))
	}
}

func (s *syslogBackend) Close() {
	s.writer.Close()
}
