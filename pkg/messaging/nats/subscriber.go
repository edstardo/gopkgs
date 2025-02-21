package nats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	natsio "github.com/nats-io/nats.go"
)

func NewSubscriber(nc *natsio.Conn, cfg *SubscriberConfig, handler natsio.MsgHandler) *Subscriber {
	return &Subscriber{
		ID:      uuid.NewString(),
		nc:      nc,
		cfg:     cfg,
		handler: handler,
	}
}

type SubscriberConfig struct {
	Subject    string
	QueueGroup string
}

type Subscriber struct {
	ID      string
	nc      *natsio.Conn
	cfg     *SubscriberConfig
	handler func(msg *natsio.Msg)
}

func (s *Subscriber) Subscribe(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	if s.cfg.Subject == "" {
		return fmt.Errorf("invalid subscriber subject")
	}

	sub, err := s.nc.Subscribe(s.cfg.Subject, s.handler)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s terminated\n", s.ID)
			return nil
		default:
			<-time.After(500 * time.Millisecond)
		}
	}
}

func (s *Subscriber) QueueSubscribe(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	if s.cfg.Subject == "" {
		return fmt.Errorf("invalid queue subscriber subject")
	}
	if s.cfg.QueueGroup == "" {
		return fmt.Errorf("invalid queue subscriber subject")
	}

	sub, err := s.nc.QueueSubscribe(s.cfg.Subject, s.cfg.QueueGroup, s.handler)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s terminated\n", s.ID)
			return nil
		default:
			<-time.After(500 * time.Millisecond)
		}
	}
}

func (s *Subscriber) SubscribeSync(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	if s.cfg.Subject == "" {
		return fmt.Errorf("invalid queue subscriber subject")
	}

	sub, err := s.nc.SubscribeSync(s.cfg.Subject)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s terminated\n", s.ID)
			return nil
		default:
			msg, err := sub.NextMsg(5 * time.Second)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
				return nil
			} else {
				s.handler(msg)
			}
		}
	}
}

func (s *Subscriber) QueueSubscribeSync(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	if s.cfg.Subject == "" {
		return fmt.Errorf("invalid queue subscriber subject")
	}
	if s.cfg.QueueGroup == "" {
		return fmt.Errorf("invalid queue subscriber subject")
	}

	sub, err := s.nc.QueueSubscribeSync(s.cfg.Subject, s.cfg.QueueGroup)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s terminated\n", s.ID)
			return nil
		default:
			msg, err := sub.NextMsg(5 * time.Second)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
				return nil
			} else {
				s.handler(msg)
			}
		}
	}
}

func (s *Subscriber) Close() {
	s.nc.Close()
}
