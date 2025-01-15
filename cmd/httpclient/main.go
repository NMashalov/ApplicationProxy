package main

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New().SetBaseURL("http://localhost:4000")
	for range time.NewTicker(5 * time.Second).C {
		response, err := client.R().Get("/ping")
		if err != nil {
			log.Print(err)
		} else {
			log.Print(string(response.Body()))
		}
	}
}
