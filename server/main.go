package main

import (
	"log"
	"log/slog"

	"github.com/franklange/go-heartbeat/utils"
)

func foo(a any) {
	log.Println()
}

func main() {
	log.SetFlags(log.Lshortfile)

	if utils.Debug() {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	foo("foo")
	foo([]string{"a", "b"})
}
