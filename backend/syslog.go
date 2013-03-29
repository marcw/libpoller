package backend

import (
	"fmt"
	"github.com/marcw/poller"
	"log/syslog"
	"os"
)

type SyslogBackend struct {
	writer *syslog.Writer
}

// Create a new SyslogBackend instance.
// Uses these environment variable:
//   - SYSLOG_NETWORK (Optionnal): "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix" and "unixpacket"."
//   - SYSLOG_ADDRESS (Optionnal): Address of the SYSLOG service
//   - SYSLOG_PREFIX (Optionnal): Prefix that will be added to the log. Defaults to "poller"
func NewSyslogBackend() (*SyslogBackend, error) {
	network := os.Getenv("SYSLOG_NETWORK")
	raddr := os.Getenv("SYSLOG_ADDRESS")
	prefix := os.Getenv("SYSLOG_PREFIX")
	if prefix == "" {
		prefix = "poller"
	}

	writer, err := syslog.Dial(network, raddr, syslog.LOG_INFO, prefix)
	if err != nil {
		return nil, err
	}

	return &SyslogBackend{writer: writer}, nil
}

func (s *SyslogBackend) Log(e *poller.Event) {
	if e.IsUp() {
		s.writer.Info(fmt.Sprintln(e.Check.Key, btos(e.IsUp()), e.Duration))
	} else {
		s.writer.Err(fmt.Sprintln(e.Check.Key, btos(e.IsUp()), e.Duration))
	}
}

func (s *SyslogBackend) Close() {
	s.writer.Close()
}
