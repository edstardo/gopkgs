package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-exit
		cancel()
	}()

	natsServer := "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224"

	nc, _ := nats.Connect(natsServer)
	defer nc.Drain()

	messages := 1_000_000
	subscribers := 1

	var wg sync.WaitGroup

	// this wg is to make sure subs are setup before pub since before
	// able to request in request-reply, existing subs are needed
	var setSubWG sync.WaitGroup

	for i := 1; i <= subscribers; i++ {
		wg.Add(1)
		setSubWG.Add(1)

		go func(id int) {
			defer wg.Done()

			var received int
			sub, err := nc.Subscribe("order.received", func(msg *nats.Msg) {
				logrus.Infof("message received: %s", string(msg.Data))
				if msg.Reply != "" {
					if err := msg.Respond(nil); err != nil {
						logrus.Error(err)
						return
					}
				}
				received++
			})

			if err != nil {
				logrus.Fatal(err)
			}
			defer sub.Unsubscribe()

			setSubWG.Done()
			<-ctx.Done()

			logrus.Infof("subscriber[%d] closed: %d reveived = %.2f percent", id, received, float64(received)/float64(messages)*100.0)
		}(i)
	}

	setSubWG.Wait() // wait for all subs to be setup

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		sentAndReceived := 0
		for sentAndReceived != messages {
			if ctx.Err() != nil {
				break
			}

			// // using ordinary publish block
			// // try using this block instead of the request-reply below
			// // to set that subscribers will not be ablt to read all messages
			// err := nc.Publish("order.received", []byte(uuid.NewString()))
			// if err != nil {
			// 	logrus.Fatal(err)
			// 	continue
			// }

			// using request-reply block
			_, err := nc.Request("order.received", []byte(uuid.NewString()), 5*time.Second)
			if err != nil {
				logrus.Fatal(err)
				continue
			}

			sentAndReceived++
		}

		<-ctx.Done()
	}(ctx)
	wg.Wait()

	// delay to make sure the graceful shutdown message is visible when the program terminates
	<-time.After(1 * time.Second)
	logrus.Info("graceful shutdown")
}
