package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	natsio "github.com/nats-io/nats.go"
)

var (
	URL         = natsio.DefaultURL
	NUM_SUBS    = 5
	NUM_MSGS    = 1_000
	QUEUE_GROUP = "test"
	SUBJECT     = "orders"
)

func main() {
	nc, err := natsio.Connect(URL)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully connect to nats server")

	var wg sync.WaitGroup

	for i := 1; i <= NUM_SUBS; i++ {
		wg.Add(1)
		go subscribe(nc, SUBJECT, i, &wg)
	}

	<-time.After(1 * time.Second)

	for i := 1; i <= NUM_MSGS; i++ {
		<-time.After(5 * time.Millisecond)
		nc.Publish(SUBJECT, []byte(fmt.Sprintf("msg-%d", i)))
	}

	wg.Wait()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
}

func subscribe(nc *natsio.Conn, subject string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	readMsgs := 0

	// nats QueueSubscribeSync example
	sub, err := nc.QueueSubscribeSync(subject, QUEUE_GROUP)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	for {
		msg, err := sub.NextMsg(5 * time.Second)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		readMsgs++
		fmt.Printf("subs-%d: %s\n", id, string(msg.Data))
	}

	// // nats SubscribeSync example
	// sub, err := nc.SubscribeSync(subject)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Unsubscribe()
	// for {
	// 	msg, err := sub.NextMsg(5 * time.Second)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		break
	// 	}
	// 	readMsgs++
	// 	fmt.Printf("subs-%d: %s\n", id, string(msg.Data))
	// }

	// nats QueueSubscribe example
	// sub, err := nc.QueueSubscribe(subject, QUEUE_GROUP, func(msg *natsio.Msg) {
	// 	readMsgs++
	// 	fmt.Printf("subs-%d: %s\n", id, string(msg.Data))
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Unsubscribe()
	// <-time.After(10 * time.Second)

	// nats Subscribe example
	// sub, err := nc.Subscribe(subject, func(msg *natsio.Msg) {
	// 	readMsgs++
	// 	fmt.Printf("subs-%d: %s\n", id, string(msg.Data))
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Unsubscribe()
	// <-time.After(10 * time.Second)

	fmt.Printf("sub-%d received %d msgs\n", id, readMsgs)
}
