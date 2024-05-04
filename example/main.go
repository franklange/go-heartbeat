package main

import (
	"fmt"
	"time"

	hb "github.com/franklange/go-heartbeat"
)

func main() {
	h := hb.Heartbeats{Timeout: 1 * time.Second}
	fmt.Println(&h)
}
