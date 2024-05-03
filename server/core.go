package main

import (
	"fmt"
	"log/slog"
	"time"
)

type Client = string

type Core struct {
	clients    map[Client][]time.Time
	actions    chan Action
	registered chan Client
	expired    chan []Client
}

func NewCore(inbuf int, outbuf int) Core {
	return Core{make(map[Client][]time.Time), make(chan Action, inbuf), make(chan string, outbuf), make(chan []string, outbuf)}
}

func (core *Core) runOne() {
	a := <-core.actions
	slog.Debug(fmt.Sprint(a))
	core.route(a)
}

func (core *Core) runAll() {
	for len(core.actions) > 0 {
		core.runOne()
	}
}

func (core *Core) runForever() {
	for {
		core.runOne()
	}
}

func (core *Core) route(a Action) {
	if a.tag == TagRegister {
		core.register(a.action.(Register))
		return
	}

	if a.tag == TagBeat {
		core.beat(a.action.(Beat))
		return
	}

	if a.tag == TagPrune {
		core.prune(a.action.(Prune))
		return
	}
	slog.Debug("unknown action")
}

func (core *Core) register(r Register) {
	_, found := core.clients[r.id]
	if found {
		slog.Debug("client exists")
		r.reply <- false
		return
	}
	core.clients[r.id] = []time.Time{}
	slog.Debug("register", "client", r.id)

	r.reply <- true
	core.registered <- r.id
}

func (core *Core) beat(b Beat) {
	ts, found := core.clients[b.id]
	if !found {
		slog.Debug("client not found")
		b.reply <- false
		return
	}

	if len(ts) == 0 {
		core.clients[b.id] = append(ts, b.ts)
		b.reply <- true
		slog.Debug("beat", "client", b.id)
		return
	}

	last := ts[len(ts)-1]
	if b.ts.Before(last) {
		slog.Debug("outdated timestamp")
		b.reply <- false
		return
	}

	if len(ts) < 5 {
		core.clients[b.id] = append(ts, b.ts)
	} else {
		core.clients[b.id] = append(ts[1:], b.ts)
	}
	slog.Debug("beat", "client", b.id)

	b.reply <- true
}

func (core *Core) prune(p Prune) {
	res := []Client{}

	for k, v := range core.clients {
		if len(v) == 0 {
			res = append(res, k)
			delete(core.clients, k)
			continue
		}

		last := v[len(v)-1]
		if p.ts.Before(last) {
			slog.Debug("prune ts in the past")
			continue
		}

		if p.ts.Sub(last) < (5 * time.Second) {
			continue
		}

		res = append(res, k)
		delete(core.clients, k)
	}
	p.reply <- res
	core.expired <- res
}
