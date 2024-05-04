package heartbeat

import (
	"sync"
	"time"
)

type Heartbeats struct {
	Timeout time.Duration

	mu  sync.Mutex
	ids map[string]time.Time
}

func NewHeartbeats(d time.Duration) *Heartbeats {
	return &Heartbeats{ids: make(map[string]time.Time), Timeout: d}
}

func (h *Heartbeats) NumClients() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return len(h.ids)
}

func (h *Heartbeats) HasClient(id string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, found := h.ids[id]
	return found
}

func (h *Heartbeats) Register(id string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.register_at(id, time.Now())
}

func (h *Heartbeats) Beat(id string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.beat_at(id, time.Now())
}

func (h *Heartbeats) Prune() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.prune_at(time.Now())
}

func (h *Heartbeats) register_at(id string, t time.Time) bool {
	_, found := h.ids[id]
	if found {
		return false
	}
	h.ids[id] = t
	return true
}

func (h *Heartbeats) beat_at(id string, t time.Time) bool {
	last, found := h.ids[id]
	if !found || t.Before(last) {
		return false
	}

	h.ids[id] = t
	return true
}

func (h *Heartbeats) prune_at(t time.Time) []string {
	var res []string
	for k, v := range h.ids {
		if t.Before(v) || (t.Sub(v) < h.Timeout) {
			continue
		}
		res = append(res, k)
	}

	return res
}
