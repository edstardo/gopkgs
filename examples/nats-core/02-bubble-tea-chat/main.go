package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/edstardo/gopkgs/pkg/messaging/nats"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	natsio "github.com/nats-io/nats.go"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[1]))
	if err != nil {
		panic(err)
	}
	environmentPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(environmentPath); err != nil {
		panic(err)
	}

	natsServer := os.Getenv("NATS_CLUSTER")
	natsSubject := os.Getenv("NATS_SUBJECT")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-exit
		cancel()
	}()

	nc, err := natsio.Connect(natsServer)
	if err != nil {
		panic(err)
	}

	engine := chatEngine{
		id:      uuid.NewString(),
		msgChan: make(chan chatMessage),
	}

	engine.subscriber = nats.NewSubscriber(nc, &nats.SubscriberConfig{
		Subject: natsSubject,
	}, engine.handleChatMessage)

	engine.publisher = nats.NewPublisher(nc, natsSubject)

	var wg sync.WaitGroup
	wg.Add(1)
	go engine.subscriber.Subscribe(ctx, &wg)

	p := tea.NewProgram(newChatModel(engine))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
