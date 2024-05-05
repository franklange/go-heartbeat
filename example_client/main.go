package main

import (
	"time"

	hb "github.com/franklange/go-heartbeat"
	"github.com/franklange/go-scanln"
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

	input := scanln.NewScanln()
	defer input.Stop()

	for {
		q := <-input.C
		if q == "q" {
			return
		}
	}
}
