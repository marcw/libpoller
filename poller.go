package main

import (
	_ "fmt"
	_ "github.com/peterbourgon/g2s"
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
	config := &Configuration{}
	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	config.Load(buffer)

	pollChannel := make(chan Check)

	go func() {
		client := &http.Client{}
		for {
			check := <-pollChannel
			go func() {
				for {
					time.Sleep(check.Interval)
					statusCode, duration, err := check.Poll(client)
					if err != nil {
						log.Println(err)
					}

					for _, v := range config.Backends {
						v.Log(&check, statusCode, duration)
					}
				}
			}()
		}
	}()

	for _, v := range config.Checks {
		pollChannel <- v
	}

	select {}
}
