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
	"reflect"

	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/mapper"
	"github.com/caiomarcatti12/nanogo/pkg/validator"
)

type InMemoryBroker struct {
	handlers []EventConsumer // Mudança para um slice de EventHandler
	logger   log.ILog
	i18n     i18n.I18N
}

type InMemoryConsumer struct {
	handlers map[string][]interface{} // Mudança para um slice de EventHandler
	logger   log.ILog
	i18n     i18n.I18N
}

func NewInMemoryBroker(logger log.ILog, i18n i18n.I18N) IEventDispatcher {
	return &InMemoryBroker{
		handlers: []EventConsumer{},
		logger:   logger,
		i18n:     i18n,
	}
}

func (i *InMemoryBroker) RegisterConsumer(eventConsumer EventConsumer) {
	// Adiciona o handler ao slice para o eventKey específico

	i.handlers = append(i.handlers, eventConsumer)
}

func (i *InMemoryBroker) Dispatch(event Event) {
	for _, consumer := range i.handlers {
		if consumer.Channel == event.Channel && consumer.Key == event.Key {
			handler, err := di.GetInstance().GetByFactory(consumer.IHandler)

			if err != nil {
				return
			}

			handlerValue := reflect.ValueOf(handler)
			// handlerType := handlerValue.Type()

			// var structName string
			// if handlerType.Kind() == reflect.Ptr {
			// 	structName = handlerType.Elem().Name()
			// } else {
			// 	structName = handlerType.Name()
			// }

			// span := ws.telemetry.StartChildSpan(structName + "::" + route.HandlerFunc)
			// defer (func() { ws.telemetry.EndSpan(span, err) })()

			method := handlerValue.MethodByName(consumer.HandlerFunc)

			if !method.IsValid() {
				// return nil, errors.InternalServerError(ws.i18n.Get("webserver.method_not_found", map[string]interface{}{"method": route.HandlerFunc, "path": route.Path}))
			}

			methodType := method.Type()
			numArgs := methodType.NumIn()
			args := make([]reflect.Value, numArgs)

			for i := 0; i < numArgs; i++ {
				paramType := methodType.In(i)

				ptrToStruct := reflect.New(paramType)
				err := mapper.InjectData(event.Data, ptrToStruct.Interface())

				if err != nil {
					// return nil, errors.InternalServerError(ws.i18n.Get("webserver.error_injecting_data", map[string]interface{}{"error": err}))
				}

				errorValidateStruct := validator.ValidateStruct(ptrToStruct.Interface())

				if errorValidateStruct != nil {
					// return nil, errorValidateStruct
				}

				args[i] = ptrToStruct.Elem()
			}

			go method.Call(args)
			// results := method.Call(args)
		}
	}
}
