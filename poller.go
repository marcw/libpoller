package main

import (
	"flag"
	"fmt"
	"github.com/marcw/poller/backend"
	"github.com/marcw/poller/check"
	"github.com/marcw/poller/poll"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	configFile = flag.String("config", "", "Configuration file to load checks from")
	httpAddr   = flag.String("hostname", "127.0.0.1:6000", "Adress/host on which the http server will listen.")
	userAgent  = flag.String("ua", "Poller (https://github.com/marcw/poller)", "User agent used by Poller")
	timeout    = flag.Duration("timeout", 10*time.Second, "Timeout")
	backends   = flag.String("backends", "stdout", "Backends to enable. Comma separated.")
)

func init() {
	flag.Parse()
}

func main() {
	if *configFile != "" {
		loadChecksFromFile()
	}
	instanciateBackends()

	poll := poll.NewHttpPoller(*userAgent, *timeout)

	for _, v := range check.Checks {
		go func(chck check.Check) {
			for {
				time.Sleep(chck.Interval)
				go func(c *check.Check) {
					event := poll.Poll(c)
					if event.Timeout {
						backend.Backends.LogTimeout(c)
					} else if event.Up {
						backend.Backends.LogSuccess(c, event.StatusCode, event.Duration)
					} else {
						backend.Backends.LogError(c, event.StatusCode, event.Duration)
					}
				}(&chck)
			}
		}(v)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/checks", check.Checks)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatalln("Unable to start http server", err)
	}

	select {}
}

func loadChecksFromFile() {
	buffer, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = check.FromJson(buffer)
	if err != nil {
		log.Fatal(err)
	}
}

func instanciateBackends() error {
	for _, v := range strings.Split(*backends, ",") {
		switch {
		case v == "statsd":
			statsd, err := backend.NewStatsdBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Statsd backend: %s", err)
			}
			backend.Add(statsd)

		case v == "stdout":
			stdout := backend.NewStdoutBackend()
			backend.Add(stdout)

		case v == "librato":
			librato, err := backend.NewLibratoBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the Librato backend: %s", err)
			}
			backend.Add(librato)

		case v == "syslog":
			syslog, err := backend.NewSyslogBackend()
			if err != nil {
				return fmt.Errorf("Impossible to instanciate the syslog backend: %s", err)
			}
			backend.Add(syslog)
		}
	}

	return nil
}
