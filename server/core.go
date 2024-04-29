package main

import (
	"fmt"
	"log/slog"
	"time"
)

type Client = string
type Timestamp = time.Time
type Timestamps = []time.Time

type Core struct {
	clients map[Client]Timestamps
}

func NewCore() Core {
	return Core{make(map[string][]time.Time)}
}

func (core *Core) add(c Client) bool {
	_, found := core.clients[c]
	if found {
		slog.Debug(fmt.Sprintln("client exists: ", c))
		return false
	}

	core.clients[c] = Timestamps{}
	return true
}

func (core *Core) beat(c Client) bool {
	return core.beat_at(c, time.Now())
}

func (core *Core) beat_at(c Client, t Timestamp) bool {
	ts, found := core.clients[c]
	if !found {
		slog.Debug("not found", "client", c)
		return false
	}

	if len(ts) == 0 {
		core.clients[c] = append(ts, t)
		return true
	}

	last := ts[len(ts)-1]
	if t.Before(last) {
		slog.Debug("outdated timestamp")
		return false
	}

	if len(ts) < 5 {
		core.clients[c] = append(ts, t)
	} else {
		core.clients[c] = append(ts[1:], t)
	}

	return true
}

func (core *Core) prune() []Client {
	return core.prune_at(time.Now())
}

func (core *Core) prune_at(t time.Time) []Client {
	res := []Client{}

	for k, v := range core.clients {
		if len(v) == 0 {
			res = append(res, k)
			delete(core.clients, k)
			continue
		}

		last := v[len(v)-1]
		if t.Before(last) {
			slog.Debug("prune ts in the past")
			continue
		}

		if t.Sub(last) < (5 * time.Second) {
			continue
		}

		res = append(res, k)
		delete(core.clients, k)
	}
	return res
}
