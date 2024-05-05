package heartbeat

import "time"

type HttpConfig struct {
	Port uint16
	Url  string
}

type GrpcConfig struct {
	Port uint16
}

type ServerConfig struct {
	HttpConfig     *HttpConfig
	GrpcConfig     *GrpcConfig
	PrunteInterval time.Duration
}

type Server struct {
	Alive chan string
	Dead  chan []string

	quit         chan bool
	heartbeats   *Heartbeats
	httpListener *httpListener
}

func NewServer(config *ServerConfig) *Server {
	s := &Server{
		Alive: make(chan string, 10),
		Dead:  make(chan []string, 10),
		quit:  make(chan bool, 1),
	}
	s.heartbeats = NewHeartbeats()
	go s.pruneRunner(config.PrunteInterval)

	if config.HttpConfig != nil {
		s.httpListener = newHttpServer(&httpListenerConfig{
			port:  config.HttpConfig.Port,
			url:   config.HttpConfig.Url,
			hbs:   s.heartbeats,
			alive: s.Alive,
		})
		go s.httpListener.run()
	}
	return s
}

func (s *Server) Stop() {
	s.quit <- true
	s.httpListener.stop()
}

func (s *Server) pruneRunner(d time.Duration) {
	t := time.NewTicker(d)
	for {
		select {
		case <-t.C:
			s.prune()
		case <-s.quit:
			return
		}
	}
}

func (s *Server) prune() {
	deads := s.heartbeats.Prune()
	if len(deads) > 0 {
		s.Dead <- deads
	}
}
