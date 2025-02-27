package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	natsio "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	NATS_SERVER     = "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224"
	NUM_PUBLISHERS  = 10
	NUM_SUBSCRIBERS = 10
)

var (
	requestRepo  = &msgRepo{}
	responseRepo = &msgRepo{}
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel()
		timer := time.NewTimer(10 * time.Second)

		for {
			select {
			case <-exit:
				return
			case <-timer.C:
				return
			}
		}
	}()

	nc, err := natsio.Connect(NATS_SERVER)
	if err != nil {
		logrus.Fatal(err)
	}
	defer nc.Drain()

	var wg sync.WaitGroup

	wg.Add(NUM_SUBSCRIBERS)
	go runResponders(mainCtx, nc, NUM_SUBSCRIBERS, &wg)

	// added delay to make sure that all responders all already set
	<-time.After(1 * time.Second)

	wg.Add(NUM_PUBLISHERS)
	go runPublishers(mainCtx, nc, NUM_PUBLISHERS, &wg)

	wg.Wait()

	<-time.After(3 * time.Second)
	logrus.Infof("num of total requests received by all responders: %d", len(requestRepo.msgs))
	logrus.Infof("num of total responses received by all publishers: %d", len(responseRepo.msgs))
}

type msgRepo struct {
	lock sync.Mutex
	msgs []string
}

func runPublishers(ctx context.Context, nc *natsio.Conn, np int, wg *sync.WaitGroup) {
	subject := "messages"

	for i := range np {
		pub := &publisher{
			nc:      nc,
			subject: subject,

			id: i + 1,
		}

		go pub.Start(ctx, wg)
	}
}

func (p *publisher) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer p.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			response, err := p.nc.Request(p.subject, fmt.Appendf(nil, "publisher[%d]-request-%d", p.id, p.sent+1), 5*time.Second)
			if err != nil {
				logrus.Errorf("error exncountered, publisher exiting: %s", err.Error())
				return
			}

			data := string(response.Data)

			responseRepo.lock.Lock()
			responseRepo.msgs = append(responseRepo.msgs, data)
			responseRepo.lock.Unlock()

			logrus.Infof("publisher[%d] received: %s", p.id, data)

			p.sent++
			p.responded++
		}
	}
}

func (p *publisher) Close() {
	<-time.After(2 * time.Second)
	logrus.Infof("publisher[%d] sent[%d] responded[%d]", p.id, p.sent, p.responded)
}

func runResponders(ctx context.Context, nc *natsio.Conn, nr int, wg *sync.WaitGroup) {
	subject := "messages"
	queueGroup := "queue"

	for id := range nr {
		resp := &responder{
			nc: nc,

			subject:    subject,
			queueGroup: queueGroup,

			id: id,
		}
		go resp.Start(ctx, wg)
	}
}

type publisher struct {
	nc      *natsio.Conn
	subject string

	id        int
	sent      int
	responded int
}

type responder struct {
	nc *natsio.Conn

	subject    string
	queueGroup string

	id        int
	received  int
	processed int
}

func (r *responder) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer r.Close()

	sub, err := r.nc.QueueSubscribe(r.subject, r.queueGroup, func(msg *natsio.Msg) {
		r.received++
		data := string(msg.Data)

		requestRepo.lock.Lock()
		requestRepo.msgs = append(requestRepo.msgs, data)
		requestRepo.lock.Unlock()

		if msg.Reply != "" {
			if err := msg.Respond(fmt.Appendf(nil, "responder[%d]-response", r.id)); err != nil {
				logrus.Fatal(err)
				return
			}
		}

		r.processed++
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			<-time.After(1 * time.Second)
		}
	}
}

func (r *responder) Close() {
	<-time.After(1 * time.Second)
	logrus.Infof("responder[%d] received[%d] processed[%d]", r.id, r.received, r.processed)
}
