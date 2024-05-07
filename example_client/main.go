package main

import (
	"time"

	hb "github.com/franklange/go-heartbeat"
)

func main() {
	c := hb.NewClient(&hb.ClientConfig{
		Id:       "c1",
		Addr:     "localhost:9000",
		Interval: 1 * time.Second,
	})
	defer c.Stop()

	for {
		time.Sleep(time.Second)
	}
}
