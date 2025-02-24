package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/edstardo/gopkgs/pkg/messaging/nats"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

	nc, err := nats.NewNatsClient(natsServer)
	if err != nil {
		log.Fatal(err)
	}

	engine := chatEngine{
		id:      uuid.NewString(),
		nc:      nc,
		subject: "chat.bubble",
		msgChan: make(chan chatMessage),
	}

	_, err = engine.nc.Subscribe(engine.subject, engine.handleChatMessage)
	if err != nil {
		log.Fatal(err)
	}

	chatModel := newChatModel(engine)

	p := tea.NewProgram(chatModel)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
