package main

import (
	"flag"
	"fmt"
	log "github.com/marcw/gogol"
	"github.com/marcw/poller"
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

type alertPool map[poller.Alerter]bool
type backendPool map[poller.Backend]bool

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

	httpProbe := poller.NewHttpProbe(*userAgent, *timeout)
	poller := poller.NewDirectPoller()

	go httpInput(config)
	go poller.Run(config.Scheduler(), backendPool, httpProbe, alerterPool)
	go config.Scheduler().Start()

	select {}
}

func instanciateBackendPool() (*backendPool, error) {
	pool := make(backendPool)
	for _, v := range strings.Split(*backends, ",") {
		switch {
		case v == "statsd":
			statsd, err := poller.NewStatsdBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the Statsd backend: %s", err)
			}
			pool[statsd] = true

		case v == "stdout":
			stdout := poller.NewStdoutBackend()
			pool[stdout] = true

		case v == "librato":
			librato, err := poller.NewLibratoBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the Librato backend: %s", err)
			}
			pool[librato] = true

		case v == "syslog":
			syslog, err := poller.NewSyslogBackend()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the syslog backend: %s", err)
			}
			pool[syslog] = true
		}
	}

	return &pool, nil
}

func instanciateAlerterPool() (*alertPool, error) {
	pool := make(alertPool)
	for _, v := range strings.Split(*alerts, ",") {
		switch {
		case v == "smtp":
			smtp, err := poller.NewSmtpAlerter()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the smtp alerter: %s", err)
			}
			pool[smtp] = true
			break
		case v == "pagerduty":
			smtp, err := poller.NewPagerDutyAlerter()
			if err != nil {
				return nil, fmt.Errorf("Impossible to instanciate the pagerduty alerter: %s", err)
			}
			pool[smtp] = true
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

func (p backendPool) Log(event *poller.Event) {
	for k, _ := range p {
		k.Log(event)
	}
}

func (p backendPool) Close() {
	for k, _ := range p {
		k.Close()
	}
}

func (p alertPool) Alert(event *poller.Event) {
	for k := range p {
		k.Alert(event)
	}
}
