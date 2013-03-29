package main

import (
	"flag"
	"fmt"
	"github.com/marcw/poller"
	"github.com/marcw/poller/alert"
	"github.com/marcw/poller/backend"
	"github.com/marcw/poller/poll"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

var (
	configFile = flag.String("config", "", "Configuration file to load checks from")
	httpAddr   = flag.String("hostname", "127.0.0.1:6060", "Adress/host on which the http server will listen.")
	userAgent  = flag.String("ua", "Poller (https://github.com/marcw/poller)", "User agent used by Poller")
	timeout    = flag.Duration("timeout", 10*time.Second, "Timeout")
	backends   = flag.String("backends", "stdout", "Backends to enable. Comma separated.")
	alerts     = flag.String("alerts", "", "Alerts to enable. Comma separated. ie: \"smtp\"")
)

func init() {
	flag.Parse()
}

func main() {
	scheduler := check.NewScheduler()
	if *configFile != "" {
		if err := loadChecksFromFile(scheduler); err != nil {
			log.Fatalln(err)
		}
	}

	backendPool, err := instanciateBackendPool()
	if err != nil {
		log.Fatalln(err)
	}
	alerterPool, err := instanciateAlerterPool()
	if err != nil {
		log.Fatalln(err)
	}

	pollerPool := poll.NewHttpPoller(*userAgent, *timeout)

	go httpInput(scheduler)
	go func(toPoll <-chan *check.Check, bp *backend.Pool, pp poll.Poller, ap *alert.Pool) {
		for {
			toPollCheck := <-toPoll
			go func(c *check.Check, b *backend.Pool, p poll.Poller, a *alert.Pool) {
				event := p.Poll(c)
				b.Log(event)
				if event.Check.ShouldAlert() {
					a.Alert(event)
				}
			}(toPollCheck, bp, pp, ap)
		}
	}(scheduler.ToPoll, backendPool, pollerPool, alerterPool)
	go scheduler.Run()

	select {}
}

func loadChecksFromFile(s *check.Scheduler) error {
	buffer, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}

	if err := s.AddFromJSON(buffer); err != nil {
		return err
	}

	return nil
}

func instanciateBackendPool() (*backend.Pool, error) {
	pool := make(backend.Pool)
	for _, v := range strings.Split(*backends, ",") {
		switch {
		case v == "statsd":
			statsd, err := backend.NewStatsdBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the Statsd backend: %s", err)
			}
			pool.Add(statsd)

		case v == "stdout":
			stdout := backend.NewStdoutBackend()
			pool.Add(stdout)

		case v == "librato":
			librato, err := backend.NewLibratoBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the Librato backend: %s", err)
			}
			pool.Add(librato)

		case v == "syslog":
			syslog, err := backend.NewSyslogBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the syslog backend: %s", err)
			}
			pool.Add(syslog)
		}
	}

	return &pool, nil
}

func instanciateAlerterPool() (*alert.Pool, error) {
	pool := make(alert.Pool)
	for _, v := range strings.Split(*alerts, ",") {
		switch {
		case v == "smtp":
			smtp, err := alert.NewSmtpAlerter()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the smtp backend: %s", err)
			}
			pool.Add(smtp)
		}
	}

	return &pool, nil
}

func httpInput(s *check.Scheduler) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/checks", s)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalln(err)
	}
}
