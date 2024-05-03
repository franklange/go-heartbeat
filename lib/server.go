package lib

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	Registered chan string
	Expired    chan []string

	conf *Config
	done chan bool
	core *Core
	http *http.Server
}

type Config struct {
	Port      uint16
	RegRoute  string
	BeatRoute string

	InBufSize  int
	OutBufSize int
}

func NewServer(conf *Config) *Server {
	var s Server
	s.conf = conf
	s.done = make(chan bool, 1)

	s.initCore()
	s.initHttp()

	go s.runCore()
	go s.runHttp()
	go s.runPrune()

	return &s
}

func (s *Server) initCore() {
	s.core = NewCore()
	s.Registered = make(chan string, 10)
	s.Expired = make(chan []string, 10)
}

func (s *Server) runCore() {
	t := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-s.done:
			return
		case <-t.C:
			s.core.runAll()
		}
	}
}

func (s *Server) runPrune() {
	t := time.NewTicker(10 * time.Second)
	deads := make(chan []string, s.conf.OutBufSize)

	for {
		select {
		case <-s.done:
			return
		case <-t.C:
			s.core.actions <- newPrune(time.Now(), deads)
			exp := <-deads
			if len(exp) == 0 {
				slog.Debug("no prune")
				continue
			}
			s.Expired <- exp
		}
	}
}

func (s *Server) initHttp() {
	s.http = &http.Server{Addr: fmt.Sprintf(":%d", s.conf.Port)}

	regHandler := RegisterHandler{0, s.core.actions}
	beatHandler := BeatHandler{0, s.core.actions}

	http.Handle(s.conf.RegRoute, &regHandler)
	http.Handle(s.conf.BeatRoute, &beatHandler)
}

func (s *Server) runHttp() {
	err := s.http.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *Server) Stop() {
	s.done <- true
	s.http.Shutdown(context.TODO())
	slog.Debug("stop")
}
