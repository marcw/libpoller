package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func TimeoutDial(timeout time.Duration) func(netw, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, timeout)
		if err != nil {
			return nil, err
		}

		return c, nil
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide, as a argument to this software, the path to the valid json configuration file")
	}
	config := &Configuration{}
	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	config.Load(buffer)

	client := &http.Client{Transport: &http.Transport{Dial: TimeoutDial(config.Timeout)}}
	for _, v := range config.Checks {
		go func(check Check) {
			for {
				time.Sleep(check.Interval)
				statusCode, duration, err := check.Poll(client)
				for _, v := range config.Backends {
					if err != nil {
						v.LogTimeout(&check)
					} else if statusCode >= 200 && statusCode < 300 {
						v.LogSuccess(&check, statusCode, duration)
					} else {
						v.LogError(&check, statusCode, duration)
					}
				}
			}
		}(v)
	}

	if config.Url != "" {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
		err = http.ListenAndServe(config.Url, nil)
		if err != nil {
			log.Fatalln("Unable to start http server", err)
		}
	}

	select {}
}
