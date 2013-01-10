package main

import (
	"github.com/marcw/poller/check"
	"github.com/marcw/poller/config"
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

	for _, v := range cfg.Checks {
		go func(chck check.Check) {
			for {
				time.Sleep(chck.Interval)
				statusCode, duration, err := chck.Poll()
				for _, v := range cfg.Backends {
					if err != nil {
						v.LogTimeout(&chck)
					} else if statusCode >= 200 && statusCode < 300 {
						v.LogSuccess(&chck, statusCode, duration)
					} else {
						v.LogError(&chck, statusCode, duration)
					}
				}
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
