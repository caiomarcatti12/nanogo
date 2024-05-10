package event

import (
	"fmt"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

type EventHandler func(Event)

type InMemoryBroker struct {
	handlers map[string]EventHandler
	logger   log.ILog
}

func NewInMemoryBroker(logger log.ILog) IEventDispatcher {
	return &InMemoryBroker{
		handlers: make(map[string]EventHandler),
		logger:   logger,
	}
}

func (b *InMemoryBroker) RegisterConsumer(eventKey string, handler EventHandler) {
	b.handlers[eventKey] = handler
}

func (b *InMemoryBroker) Dispatch(event Event) {
	if handler, exists := b.handlers[event.Key]; exists {
		handler(event)
	} else {
		b.logger.Warning(fmt.Sprintf("Nenhum manipulador encontrado para o evento: %s", event.Key))
	}
}
