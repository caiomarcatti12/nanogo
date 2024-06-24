/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package event

import (
	"fmt"

	"github.com/caiomarcatti12/nanogo/v2/src/log"
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
