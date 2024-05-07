# go-heartbeat

Simple server and client libraries to monitor client connectivity via heartbeats using gRPC client-side streaming.  
The server library provides client connects/disconnects via channels `Alive` and `Dead`.  
The client library uses `Interval` to either automatically (re)connect to the given server address or to send heartbeats for the provided client id.  

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
			fmt.Println("conn:", c)
		case c := <-s.Dead:
			fmt.Println("dead:", c)
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