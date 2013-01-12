package main

import (
	"github.com/marcw/poller/check"
	"github.com/marcw/poller/config"
	"github.com/marcw/poller/poll"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide, as a argument to this software, the path to the valid json configuration file")
	}
	cfg := config.NewConfig()
	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	cfg.Load(buffer)

	poll := poll.NewHttpPoller()

	for _, v := range cfg.Checks {
		go func(chck check.Check) {
			for {
				time.Sleep(chck.Interval)
				go func(c *check.Check) {
					event := poll.Poll(c)
					for _, v := range cfg.Backends {
						if event.Timeout {
							v.LogTimeout(c)
						} else if event.Up {
							v.LogSuccess(c, event.StatusCode, event.Duration)
						} else {
							v.LogError(c, event.StatusCode, event.Duration)
						}
					}
				}(&chck)
			}
		}(v)
	}

	if cfg.Url != "" {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
		err = http.ListenAndServe(cfg.Url, nil)
		if err != nil {
			log.Fatalln("Unable to start http server", err)
		}
	}

	select {}
}
