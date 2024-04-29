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
}
