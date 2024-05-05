package main

import (
	"fmt"
	"time"

	hb "github.com/franklange/go-heartbeat"
	"github.com/franklange/go-scanln"
)

func main() {
	s := hb.NewServer(&hb.ServerConfig{
		HttpConfig:     &hb.HttpConfig{Port: 9000, Url: "/heartbeats"},
		PrunteInterval: 3 * time.Second,
	})
	defer s.Stop()

	input := scanln.NewScanln()
	defer input.Stop()

	for {
		select {
		case c := <-s.Alive:
			fmt.Println("alive: ", c)
		case d := <-s.Dead:
			fmt.Println("dead: ", d)
		case i := <-input.C:
			if i == "q" {
				return
			}
		}
	}
}
