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

func main() {
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

	wg.Add(1)
	go subscriberDemo(ctx, &wg)

	wg.Add(1)
	go queueSubscriberDemo(ctx, &wg)

	wg.Wait()
}

type OrdersSubscriber struct {
	Subscriber *nats.Subscriber
	ReadMsgs   int
}

func (s *OrdersSubscriber) OrdersHandler(msg *natsio.Msg) {
	s.ReadMsgs += 1
	fmt.Printf("%s: %s\n", s.Subscriber.ID, string(msg.Data))
}

func subscriberDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	url := natsio.DefaultURL
	numSubs := 5
	numMsgs := 100_000
	subject := "orders.ph"

	nc, err := natsio.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	defer mainWG.Done()

	var wg sync.WaitGroup

	var subs []*OrdersSubscriber
	go func() {
		for i := 1; i <= numSubs; i++ {
			oSub := &OrdersSubscriber{}
			oSub.Subscriber = nats.NewSubscriber(nc, &nats.SubscriberConfig{
				Subject: subject,
			}, oSub.OrdersHandler)

			subs = append(subs, oSub)
			wg.Add(1)
			go oSub.Subscriber.Subscribe(ctx, &wg)
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(1 * time.Second)

	ordersPublisher := nats.NewPublisher(nc, subject)
	for i := 1; i <= numMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= subscriberDemo demo result")
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
	}
}

func queueSubscriberDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	url := natsio.DefaultURL
	numSubs := 5
	numMsgs := 100_000
	subject := "orders.sg.queue"
	queue := "orders.sg"

	nc, err := natsio.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	defer mainWG.Done()

	var wg sync.WaitGroup

	var subs []*OrdersSubscriber
	go func() {
		for i := 1; i <= numSubs; i++ {
			oSub := &OrdersSubscriber{}
			oSub.Subscriber = nats.NewSubscriber(nc, &nats.SubscriberConfig{
				Subject:    subject,
				QueueGroup: queue,
			}, oSub.OrdersHandler)

			subs = append(subs, oSub)
			wg.Add(1)
			go oSub.Subscriber.QueueSubscribe(ctx, &wg)
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(1 * time.Second)

	ordersPublisher := nats.NewPublisher(nc, subject)
	for i := 1; i <= numMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= queueSubscriberDemo demo result")
	total := 0
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
		total += s.ReadMsgs
	}

	fmt.Printf("%d total read messages\n", total)
}
