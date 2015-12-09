package message

import (
	"fmt"
)

type Message struct {
	Topic string
	Body  string
}

type MessageCenter struct {
	Input chan *Message
}

func NewMessageCenter() *MessageCenter {
	return &MessageCenter{
		Input: make(chan *Message, 1024),
	}
}

func (mc *MessageCenter) Run() {
	for {
		select {
		case msg := <-mc.Input:
			fmt.Printf("Topic:%s Body: %s\n", msg.Topic, msg.Body)
		}
	}
}
