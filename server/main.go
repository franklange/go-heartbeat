package main

import (
	"log"
	"log/slog"

	"github.com/franklange/go-heartbeat/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if utils.Debug() {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	hbServer := NewHeartbeatServer(HeartbeatConf{
		port:       9000,
		regRoute:   "/register",
		beatRoute:  "/beat",
		inBufSize:  10,
		outBufSize: 10,
	})

	select {
	case id := <-hbServer.registered:
		slog.Info("new client", "client", id)
	case deads := <-hbServer.expired:
		slog.Info("clients lost", "expired", deads)
	}

	hbServer.stop()
}
