package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type HeartbeatServer struct {
	core       Core
	running    atomic.Bool
	registered <-chan string
	expired    <-chan []string
	httpSrv    *http.Server
}

type HeartbeatConf struct {
	port      uint16
	regRoute  string
	beatRoute string

	inBufSize  int
	outBufSize int
}

func NewHeartbeatServer(conf HeartbeatConf) *HeartbeatServer {
	var hb HeartbeatServer

	// init and run core logic
	hb.core = NewCore(conf.inBufSize, conf.outBufSize)
	hb.running.Store(true)
	hb.registered = hb.core.registered
	hb.expired = hb.core.expired

	go func() {
		for hb.running.Load() == true {
			hb.core.runAll()
		}
	}()

	// init and run http endpoint
	hb.httpSrv = &http.Server{Addr: fmt.Sprintf(":%d", conf.port)}

	regHandler := RegisterHandler{0, hb.core.actions}
	beatHandler := BeatHandler{0, hb.core.actions}

	http.Handle(conf.regRoute, &regHandler)
	http.Handle(conf.beatRoute, &beatHandler)

	go func() {
		err := hb.httpSrv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return &hb
}

func (hb *HeartbeatServer) stop() {
	hb.running.Store(false)
	hb.httpSrv.Shutdown(context.TODO())
	slog.Debug("hb stop")
}
