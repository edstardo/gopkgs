package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/edstardo/gopkgs/pkg/messaging/nats"
	"github.com/joho/godotenv"
	natsio "github.com/nats-io/nats.go"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	environmentPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(environmentPath); err != nil {
		panic(err)
	}

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

	wg.Add(1)
	go syncSubscriberDemo(ctx, &wg)

	wg.Add(1)
	go syncQueueSubscriberDemo(ctx, &wg)

	wg.Add(1)
	go subscriberChanDemo(ctx, &wg)

	wg.Add(1)
	go queueSubscriberChanDemo(ctx, &wg)

	wg.Add(1)
	go queueSubscriberSyncWithChanDemo(ctx, &wg)

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
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 10_000
	subject := "orders.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 10_000
	subject := "orders.queue.sub.demo"
	queue := "orders.queue.sub.demo"

	nc, err := natsio.Connect(natsServer)
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

func syncSubscriberDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 1_000
	subject := "orders.sync.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
			go oSub.Subscriber.SubscribeSync(ctx, &wg)
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
		<-time.After(10 * time.Millisecond)
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= syncSubscriberDemo demo result")
	total := 0
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
		total += s.ReadMsgs
	}
}

func syncQueueSubscriberDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 1_000
	subject := "orders.sync.queue.sub.demo"
	queue := "orders.sync.queue.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
			go oSub.Subscriber.QueueSubscribeSync(ctx, &wg)
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
		<-time.After(10 * time.Millisecond)
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= syncQueueSubscriberDemo demo result")
	total := 0
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
		total += s.ReadMsgs
	}

	fmt.Printf("%d total read messages\n", total)
}

func subscriberChanDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 10
	numMsgs := 1_000
	subject := "orders.chan.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
			go oSub.Subscriber.SubscribeWithChan(ctx, &wg, make(chan *natsio.Msg))
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(2 * time.Second)

	ordersPublisher := nats.NewPublisher(nc, subject)
	for i := 1; i <= numMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		// publish msg delay
		<-time.After(10 * time.Millisecond)
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= subscriberChanDemo demo result")
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
	}
}

func queueSubscriberChanDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 1_000
	subject := "orders.chan.queue.sub.demo"
	queue := "orders.chan.queue.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
			go oSub.Subscriber.QueueSubscribeWithChan(ctx, &wg, make(chan *natsio.Msg))
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(2 * time.Second)

	ordersPublisher := nats.NewPublisher(nc, subject)
	for i := 1; i <= numMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		// publish msg delay
		<-time.After(10 * time.Millisecond)
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= queueSubscriberChanDemo demo result")
	var total int
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
		total += s.ReadMsgs
	}
	fmt.Printf("%d total read messages\n", total)
}

func queueSubscriberSyncWithChanDemo(ctx context.Context, mainWG *sync.WaitGroup) {
	natsServer := os.Getenv("NATS_CLUSTER")
	numSubs := 5
	numMsgs := 1_000
	subject := "orders.chan.queue.sync.sub.demo"
	queue := "orders.chan.queue.sync.sub.demo"

	nc, err := natsio.Connect(natsServer)
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
			go oSub.Subscriber.QueueSubscribeSyncWithChan(ctx, &wg, make(chan *natsio.Msg))
		}
	}()

	// add delay to make sure all consumers are already
	// setup and ready to receive messages
	<-time.After(2 * time.Second)

	ordersPublisher := nats.NewPublisher(nc, subject)
	for i := 1; i <= numMsgs; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		// publish msg delay
		<-time.After(10 * time.Millisecond)
		if err := ordersPublisher.Publish([]byte(fmt.Sprintf("msg-%d", i))); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()

	fmt.Println("\n============================= queueSubscriberSyncWithChanDemo demo result")
	var total int
	for _, s := range subs {
		fmt.Printf("%s read %d messages\n", s.Subscriber.ID, s.ReadMsgs)
		total += s.ReadMsgs
	}
	fmt.Printf("%d total read messages\n", total)
}
