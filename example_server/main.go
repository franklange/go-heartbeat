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
			fmt.Println("alive: ", c)
		case c := <-s.Dead:
			fmt.Println("dead: ", c)
		}
	}
}
