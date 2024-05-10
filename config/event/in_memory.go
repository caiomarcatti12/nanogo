package event

import (
	"fmt"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

type EventHandler func(Event)

type InMemoryBroker struct {
	handlers map[string][]EventHandler // Mudança para um slice de EventHandler
	logger   log.ILog
}

func NewInMemoryBroker(logger log.ILog) IEventDispatcher {
	return &InMemoryBroker{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

func (b *InMemoryBroker) RegisterConsumer(eventKey string, handler EventHandler) {
	// Adiciona o handler ao slice para o eventKey específico
	b.handlers[eventKey] = append(b.handlers[eventKey], handler)
}

func (b *InMemoryBroker) Dispatch(event Event) {
	if handlers, exists := b.handlers[event.Key]; exists {
		for _, handler := range handlers {
			go handler(event) // Executa cada handler em uma nova goroutine
		}
	} else {
		b.logger.Warning(fmt.Sprintf("Nenhum manipulador encontrado para o evento: %s", event.Key))
	}
}
