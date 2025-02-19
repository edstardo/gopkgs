package nats

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsClient{conn: nc}, nil
}

func (c *NatsClient) Close() {
	c.conn.Close()
}

func (c *NatsClient) Publish(subject string, data []byte) error {
	return c.conn.Publish(subject, data)
}

func (c *NatsClient) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}
