package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/edstardo/gopkgs/pkg/messaging/nats"
	natsio "github.com/nats-io/nats.go"
)

var (
	URL           = natsio.DefaultURL
	NumConsumers  = 5
	NumMsgs       = 500000
	QueueGroup    = "test"
	SubjectOrders = "orders"
)

func main() {
	nc, err := nats.NewNatsClient(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-exit:
				fmt.Println("exit chan triggered")
				cancel()
				return
			default:
				<-time.After(1 * time.Second)
			}
		}
	}()

	var wg sync.WaitGroup
	var consumers []*Ordersconsumer
	go func() {
		for i := 1; i <= NumConsumers; i++ {
			cons := &Ordersconsumer{
				NC:      nc,
				Subject: SubjectOrders,
				ID:      fmt.Sprintf("consumer-%d", i),
			}
			consumers = append(consumers, cons)
			wg.Add(1)
			go cons.Start(ctx, &wg)
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(1 * time.Second)

	for i := 1; i <= NumMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		if err := nc.Publish(SubjectOrders, []byte(fmt.Sprintf("msg-%d", i))); err != nil {
			panic(err)
		}
	}

	wg.Wait()

	for _, c := range consumers {
		fmt.Printf("%s read %d messages\n", c.ID, c.ReadMsgs)
	}
}

type Ordersconsumer struct {
	ID       string
	NC       *nats.NatsClient
	Subject  string
	ReadMsgs int
}

func (s *Ordersconsumer) OrdersHandler(msg *natsio.Msg) {
	s.ReadMsgs += 1
	fmt.Printf("%s: %s\n", s.ID, string(msg.Data))
}

func (s *Ordersconsumer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	sub, err := s.NC.Subscribe(s.Subject, s.OrdersHandler)
	if err != nil {
		return err
	}
	defer wg.Done()
	defer sub.Unsubscribe()

	fmt.Printf("%s started...\n", s.ID)

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
