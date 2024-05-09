package eda

import "fmt"

var _ IEventDispatcher = &InMemoryBroker{}

type EventHandler func(Event)

type InMemoryBroker struct {
	handlers map[string]EventHandler
}

func NewInMemoryBroker() IEventDispatcher {
	return &InMemoryBroker{
		handlers: make(map[string]EventHandler),
	}
}

func (b *InMemoryBroker) RegisterEvent(eventKey string, handler EventHandler) {
	b.handlers[eventKey] = handler
}

func (b *InMemoryBroker) SendEvent(event Event) {
	if handler, exists := b.handlers[event.Key]; exists {
		handler(event)
	} else {
		fmt.Println("No handler registered for event:", event.Key)
	}
}
