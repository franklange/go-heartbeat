package main

import (
	"fmt"
	"time"

	hb "github.com/franklange/go-heartbeat"
)

func main() {
	hbc := hb.NewClient(&hb.ClientConfig{
		Host:     "localhost",
		Port:     9000,
		Url:      "/heartbeats",
		ClientId: "c1",
		Interval: 1 * time.Second,
	})
	defer hbc.Stop()

	var input string
	for {
		fmt.Scanln(&input)
		if input == "q" {
			break
		}
		fmt.Println(input)
	}
}
