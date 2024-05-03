package lib

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type HttpListener struct {
	regs  RegisterHandler
	beats BeatHandler
}

type HttpConfig struct {
	port      uint16
	regRoute  string
	beatRoute string
}

type HttpRegister struct {
	Id string `json:"id"`
}

type HttpBeat struct {
	Id string
}

type RegisterHandler struct {
	nextId  uint
	actions chan<- Action
}

func decode[T any](r *http.Request) (*T, error) {
	var t T
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&t)
	if err != nil {
		slog.Warn("decode")
		return nil, err
	}

	return &t, nil
}

func (handler *RegisterHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	regHttp, err := decode[HttpRegister](request)
	if err != nil || regHttp.Id == "" {
		slog.Warn("id missing", "method", "register")
		return
	}

	reply := make(chan bool, 1)
	handler.actions <- newRegister(regHttp.Id, reply)

	ok := <-reply
	if !ok {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 - Client already registered"))
	}
}

type BeatHandler struct {
	nextId  uint
	actions chan<- Action
}

func (handler *BeatHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	beatHttp, err := decode[HttpBeat](request)
	if err != nil || beatHttp.Id == "" {
		slog.Warn("id missing", "method", "beat")
		return
	}

	reply := make(chan bool, 1)
	handler.actions <- newBeat(beatHttp.Id, time.Now(), reply)

	ok := <-reply
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Unknown client id"))
	}
}
