package heartbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type beatRequest struct {
	Id string
}

type httpListenerConfig struct {
	port  uint16
	url   string
	hbs   *Heartbeats
	alive chan<- string
}

type httpListener struct {
	httpServer *http.Server
	heartbeats *Heartbeats
	alive      chan<- string
}

func newHttpServer(config *httpListenerConfig) *httpListener {
	h := &httpListener{
		heartbeats: config.hbs,
		alive:      config.alive,
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", config.port)},
	}
	http.HandleFunc(config.url, h.postBeat())

	return h
}

func (h *httpListener) run() {
	err := h.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (h *httpListener) stop() {
	h.httpServer.Shutdown(context.TODO())
}

func (h *httpListener) postBeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var beat beatRequest
		err := json.NewDecoder(r.Body).Decode(&beat)
		if err != nil || beat.Id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("client id missing"))
			return
		}
		count := h.heartbeats.Beat(beat.Id)
		if count > 3 {
			h.alive <- beat.Id
		}
	}
}
