package mynats

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

func (c *NatsClient) Publish(subject string, message []byte) error {
    return c.Conn.Publish(subject, message)
}

func (c *NatsClient) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
    return c.Conn.Subscribe(subject, handler)
}

func (c *NatsClient) Close() {
    c.Conn.Close()
}
