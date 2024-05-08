package main

import (
	"fmt"

	"github.com/franklange/go-heartbeat"
)

func main() {
	s := heartbeat.NewHeartbeatServer("9000")
	defer s.Stop()

	for {
		select {
		case c := <-s.Alive:
			fmt.Printf("[conn] id: %s addr: %s\n", c.Id, c.Addr)
		case c := <-s.Dead:
			fmt.Printf("[dead] id: %s addr: %s\n", c.Id, c.Addr)
		}
	}
}
