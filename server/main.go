package main

import (
	"log"
	"log/slog"

	hb "github.com/franklange/go-heartbeat/lib"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	if hb.Debug() {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	hbServer := hb.NewServer(&hb.Config{
		Port:       9000,
		RegRoute:   "/register",
		BeatRoute:  "/beat",
		InBufSize:  10,
		OutBufSize: 10,
	})
	defer hbServer.Stop()

	for {
		select {
		case id := <-hbServer.Registered:
			slog.Info("new client", "client", id)
		case deads := <-hbServer.Expired:
			slog.Info("clients lost", "expired", deads)
		}
	}
}
