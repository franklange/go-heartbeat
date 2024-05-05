package heartbeat

import (
	"sync"
)

type Heartbeats struct {
	mu      sync.Mutex
	clients map[string]int
}

func NewHeartbeats() *Heartbeats {
	return &Heartbeats{clients: make(map[string]int)}
}

func (h *Heartbeats) NumClients() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return len(h.clients)
}

func (h *Heartbeats) HasClient(id string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, found := h.clients[id]
	return found
}

func (h *Heartbeats) Beat(id string) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	count, found := h.clients[id]
	if !found {
		h.clients[id] = 3
	} else {
		h.clients[id] = min((count + 1), 3)
	}
	return h.clients[id]
}

func (h *Heartbeats) Prune() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	var res []string
	for client, count := range h.clients {
		count--
		if count <= 0 {
			res = append(res, client)
			delete(h.clients, client)
		} else {
			h.clients[client] = count
		}
	}
	return res
}

func (h *Heartbeats) numBeats(id string) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	beats, found := h.clients[id]
	if !found {
		return -1
	}
	return beats
}
