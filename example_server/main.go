package main

import (
	"fmt"
	"time"

	hb "github.com/franklange/go-heartbeat"
)

func main() {
	s := hb.NewServer(&hb.ServerConfig{
		HttpConfig:     &hb.HttpConfig{Port: 9000, Url: "/heartbeats"},
		PrunteInterval: 5 * time.Second,
	})
	defer s.Stop()

	for {
		select {
		case c := <-s.Alive:
			fmt.Println("alive: ", c)
			break
		case cs := <-s.Dead:
			fmt.Println("dead: ", cs)
			break
		}
	}
}
