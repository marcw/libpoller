package main

import (
	"flag"
	"fmt"
	"github.com/marcw/poller"
	"github.com/marcw/poller/alert"
	"github.com/marcw/poller/backend"
	"github.com/marcw/poller/service"
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
	var store poller.Store

	if *configFile != "" {
		store = poller.NewFileStore(*configFile)
	} else {
		store = poller.NewInMemoryStore()
	}
	config := poller.NewConfig(store, poller.NewSimpleScheduler())
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	backendPool, err := instanciateBackendPool()
	if err != nil {
		log.Fatalln(err)
	}
	alerterPool, err := instanciateAlerterPool()
	if err != nil {
		log.Fatalln(err)
	}

	pollerPool := service.NewHttpPoller(*userAgent, *timeout)
	poller := poller.NewDirectPoller()

	go httpInput(config)
	go poller.Run(config.Scheduler(), backendPool, pollerPool, alerterPool)
	go config.Scheduler().Start()

	select {}
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
				return nil, fmt.Errorf("Impossible to instanciate the smtp alerter: %s", err)
			}
			pool.Add(smtp)
			break
		case v == "pagerduty":
			smtp, err := alert.NewPagerDutyAlerter()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the pagerduty alerter: %s", err)
			}
			pool.Add(smtp)
		}
	}

	return &pool, nil
}

func httpInput(config *poller.Config) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/checks", poller.NewConfigHttpHandler(config))
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalln(err)
	}
}
