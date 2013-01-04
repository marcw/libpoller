package main

import (
	"fmt"
	"github.com/peterbourgon/g2s"
	"github.com/rcrowley/go-librato"
	"log"
	"log/syslog"
	"os"
	"time"
)

type Backend interface {
	// "Log" check result in the Backend service
	LogSuccess(check *Check, statusCode int, duration time.Duration)
	LogError(Check *Check, statusCode int, duration time.Duration)
	LogTimeout(check *Check)
	Close()
}

// Backend for Statsd
type StatsdBackend struct {
	statsd g2s.Statter
	prefix string
}

// Instanciate a new StatsdBackend
// Uses:
//   STATSD_HOST env variable
//   STATSD_PORT env variable (defaults to 8125)
//   STATSD_PROTOCOL env variable (defaults to udp)
//   STATSD_PREFIX env variable (defaults to `poller.checks.`)
func NewStatsdBackend() (*StatsdBackend, error) {
	envHost := os.Getenv("STATSD_HOST")
	envPort := os.Getenv("STATSD_PORT")
	envProtocol := os.Getenv("STATSD_PROTOCOL")
	envPrefix := os.Getenv("STATSD_PREFIX")

	if envHost == "" {
		return nil, fmt.Errorf("STATSD_HOST environment variable must be defined")
	}

	if envPort == "" {
		envPort = "8125"
	}

	if envProtocol == "" {
		envProtocol = "udp"
	}

	if envPrefix == "" {
		envPrefix = "poller.checks."
	}

	statsd, err := g2s.Dial(envProtocol, envHost+":"+envPort)
	if err != nil {
		return nil, err
	}

	return &StatsdBackend{statsd: statsd, prefix: envPrefix}, nil
}

// Log to statsd the check result as follow:
// `Check.Key`.duration : Request duration
// `Check.Key`.success : Request succeeded (status code is 2xx)
// `Check.Key`.error : Request failed (status code != 200)
func (s *StatsdBackend) LogSuccess(check *Check, statusCode int, duration time.Duration) {
	s.logDuration(check, duration)
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 1)
}

func (s *StatsdBackend) LogError(check *Check, statusCode int, duration time.Duration) {
	s.logDuration(check, duration)
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 0)
}

func (s *StatsdBackend) LogTimeout(check *Check) {
	s.statsd.Counter(1.0, s.prefix+check.Key+".up", 0)
}

func (s *StatsdBackend) logDuration(check *Check, duration time.Duration) {
	s.statsd.Timing(1.0, s.prefix+check.Key+".duration", duration)
}

func (s *StatsdBackend) Close() {
	// NO OP
}

// StdoutBackend logs result to Stdout
type StdoutBackend struct {
}

func NewStdoutBackend() *StdoutBackend {

	return &StdoutBackend{}
}

func (s *StdoutBackend) LogSuccess(check *Check, statusCode int, duration time.Duration) {
	log.Println(check.Key, statusCode, duration)
}

func (s *StdoutBackend) LogError(check *Check, statusCode int, duration time.Duration) {
	log.Println(check.Key, statusCode, duration)
}

func (s *StdoutBackend) LogTimeout(check *Check) {
	log.Println(check.Key, "TIMEOUT")
}

func (s *StdoutBackend) Close() {
	// NO OP
}

type LibratoBackend struct {
	metrics librato.Metrics
	prefix  string
}

func NewLibratoBackend() (*LibratoBackend, error) {
	user := os.Getenv("LIBRATO_USER")
	token := os.Getenv("LIBRATO_TOKEN")
	source := os.Getenv("LIBRATO_SOURCE")
	prefix := os.Getenv("LIBRATO_PREFIX")

	if user == "" {
		return nil, fmt.Errorf("LIBRATO_USER environment variable must be defined")
	}

	if token == "" {
		return nil, fmt.Errorf("LIBRATO_TOKEN environment variable must be defined")
	}

	if source == "" {
		source = "poller"
	}

	if prefix == "" {
		prefix = "poller.checks."
	}

	metrics := librato.NewSimpleMetrics(user, token, source)

	return &LibratoBackend{metrics: metrics, prefix: prefix}, nil
}

func (l *LibratoBackend) LogSuccess(check *Check, statusCode int, duration time.Duration) {
	l.logDuration(check, duration)
	c := l.metrics.GetGauge(l.prefix + check.Key + ".up")
	c <- 1
}

func (l *LibratoBackend) LogError(check *Check, statusCode int, duration time.Duration) {
	l.logDuration(check, duration)
	c := l.metrics.GetGauge(l.prefix + check.Key + ".up")
	c <- 0
}

func (l *LibratoBackend) LogTimeout(check *Check) {
	c := l.metrics.GetGauge(l.prefix + check.Key + ".up")
	c <- 0
}

func (l *LibratoBackend) logDuration(check *Check, duration time.Duration) {
	d := l.metrics.GetGauge(l.prefix + check.Key + ".duration")
	d <- int64(duration.Nanoseconds() / int64(time.Millisecond))
}

func (l *LibratoBackend) Close() {
	l.metrics.Close()
	l.metrics.Wait()
}

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

func (s *SyslogBackend) LogSuccess(check *Check, statusCode int, duration time.Duration) {
	s.writer.Info(fmt.Sprintln(check.Key, statusCode, duration))
}

func (s *SyslogBackend) LogError(check *Check, statusCode int, duration time.Duration) {
	s.writer.Err(fmt.Sprintln(check.Key, statusCode, duration))
}

func (s *SyslogBackend) LogTimeout(check *Check) {
	s.writer.Err(fmt.Sprintln(check.Key, "TIMEOUT"))
}

func (s *SyslogBackend) Close() {
	s.writer.Close()
}
