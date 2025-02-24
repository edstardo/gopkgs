package main

import (
	"encoding/json"

	"github.com/edstardo/gopkgs/pkg/messaging/nats"
	natsio "github.com/nats-io/nats.go"
)

type chatEngine struct {
	id         string
	subscriber *nats.Subscriber
	publisher  *nats.Publisher
	msgChan    chan chatMessage
}

func (c *chatEngine) handleChatMessage(msg *natsio.Msg) {
	var chat chatMessage
	if err := json.Unmarshal(msg.Data, &chat); err != nil {
		panic(err)
	}
	c.msgChan <- chat
}

func (c *chatEngine) sendChatMessage(message string) {
	chat := chatMessage{
		SenderID: c.id,
		Message:  message,
	}

	data, err := json.Marshal(chat)
	if err != nil {
		panic(err)
	}

	if err := c.publisher.Publish(data); err != nil {
		panic(err)
	}
}

type chatMessage struct {
	SenderID string `json:"sender_id"`
	Message  string `json:"message"`
}
