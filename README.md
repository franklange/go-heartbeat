# go-heartbeat

Simple server/client library to monitor connectivity via heartbeats using gRPC client-side streaming.  
The server provides `Peer` events with `Id` and `Addr` via channels `Alive` and `Dead`.  
The client automatically (re)connects and sends heartbeats with the user-provided `Id` with `Interval` frequency.  

### Example Server

```go
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

```

### Example Client
```go
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
```