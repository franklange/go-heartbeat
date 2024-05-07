package heartbeat

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/franklange/go-heartbeat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id     string
	quit   chan bool
	conn   *grpc.ClientConn
	client proto.HeartbeatClient
	stream proto.Heartbeat_ConnectClient
}

type ClientConfig struct {
	Id       string
	Addr     string
	Interval time.Duration
}

func NewClient(config *ClientConfig) *Client {
	conn, err := grpc.NewClient(config.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}

	c := &Client{
		id:     config.Id,
		quit:   make(chan bool, 1),
		conn:   conn,
		client: proto.NewHeartbeatClient(conn),
	}
	go c.run(config.Interval)

	return c
}

func (c *Client) Stop() {
	c.quit <- true
	if c.stream != nil {
		c.stream.CloseAndRecv()
	}
	c.conn.Close()
}

func (c *Client) run(d time.Duration) {
	c.beat()
	t := time.NewTicker(d)
	for {
		select {
		case <-c.quit:
			return
		case <-t.C:
			c.beat()
		}
	}
}

func (c *Client) beat() {
	if c.stream == nil {
		stream, err := c.client.Connect(context.Background())
		if err != nil {
			return
		}
		c.stream = stream
	}

	err := c.stream.Send(&proto.Beat{ClientId: c.id})
	if err == io.EOF {
		c.stream = nil
	}
}
