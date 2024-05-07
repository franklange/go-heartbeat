package heartbeat

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/franklange/go-heartbeat/proto"
	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedHeartbeatServer
	grpcServer *grpc.Server

	Dead  chan string
	Alive chan string
}

func NewHeartbeatServer(port string) *Server {
	s := &Server{
		Dead:       make(chan string, 100),
		Alive:      make(chan string, 100),
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
	var id string
	for {
		beat, err := stream.Recv()
		if err != nil {
			if (err == io.EOF) && (id != "") {
				s.Dead <- id
				return nil
			}
			return err
		}
		if id != "" {
			continue
		}

		id = beat.ClientId
		s.Alive <- id
	}
}
