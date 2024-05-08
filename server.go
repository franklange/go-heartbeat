package heartbeat

import (
	"fmt"
	"log"
	"net"

	"github.com/franklange/go-heartbeat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Peer struct {
	Id   string
	Addr string
}

type Server struct {
	proto.UnimplementedHeartbeatServer
	grpcServer *grpc.Server

	Dead  chan Peer
	Alive chan Peer
}

func NewHeartbeatServer(port string) *Server {
	s := &Server{
		Dead:       make(chan Peer, 100),
		Alive:      make(chan Peer, 100),
		grpcServer: grpc.NewServer(),
	}

	lis, err := net.Listen("tcp", fmt.Sprint("localhost:", port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	proto.RegisterHeartbeatServer(s.grpcServer, s)

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	return s
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) Connect(stream proto.Heartbeat_ConnectServer) error {
	var p Peer
	for {
		beat, err := stream.Recv()
		if err != nil {
			if p.Id != "" {
				s.Dead <- p
			}
			return err
		}
		if p.Id != "" {
			continue
		}

		p.Id = beat.ClientId
		pp, ok := peer.FromContext(stream.Context())
		if !ok {
			p.Addr = "unknown"
		} else {
			p.Addr = pp.Addr.String()
		}

		s.Alive <- p
	}

}
