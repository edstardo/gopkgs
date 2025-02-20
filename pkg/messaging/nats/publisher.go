package nats

import (
	natsio "github.com/nats-io/nats.go"
)

type Publisher struct {
	nc      *natsio.Conn
	subject string
}

func NewPublisher(nc *natsio.Conn, subject string) *Publisher {
	return &Publisher{
		nc:      nc,
		subject: subject,
	}
}

func (p *Publisher) Publish(data []byte) error {
	return p.nc.Publish(p.subject, data)
}
