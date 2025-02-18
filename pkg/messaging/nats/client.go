package nats

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
    Conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }
    return &NatsClient{Conn: nc}, nil
}

func (c *NatsClient) Close() {
    c.Conn.Close()
}
