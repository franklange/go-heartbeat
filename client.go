package heartbeat

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type ClientConfig struct {
	Port     uint16
	Host     string
	Url      string
	ClientId string
	Interval time.Duration
}

type Client struct {
	posturl    string
	payload    []byte
	quit       chan bool
	httpClient *http.Client
	interval   time.Duration
}

func NewClient(config *ClientConfig) *Client {
	c := &Client{
		posturl:    fmt.Sprintf("http://%s:%d%s", config.Host, config.Port, config.Url),
		payload:    []byte(fmt.Sprintf("{\"id\": \"%s\"}", config.ClientId)),
		quit:       make(chan bool, 1),
		httpClient: &http.Client{},
		interval:   config.Interval,
	}
	go c.run()

	return c
}

func (c *Client) Stop() {
	c.quit <- true
}

func (c *Client) run() {
	c.send()
	t := time.NewTicker(c.interval)

	for {
		select {
		case <-c.quit:
			return
		case <-t.C:
			c.send()
		}
	}
}

func (c *Client) send() {
	request, _ := http.NewRequest("POST", c.posturl, bytes.NewBuffer(c.payload))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c.httpClient.Do(request)
}
