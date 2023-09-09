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
package rabbitmq

import (
	"encoding/json"
	"log"
)

type Consumer interface {
	Consume(body map[string]interface{}, headers map[string]interface{})
}

func Consume(exchange Exchange, queue Queue, consumer Consumer) {
	connection := NewInstanceRabbitmq()

	DeclareExchange(exchange)
	DeclareQueue(queue)

	BindQueue(exchange, queue)

	msgs, err := connection.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			bodyMap := make(map[string]interface{})

			err := json.Unmarshal(d.Body, &bodyMap)
			if err != nil {
				log.Fatalf("Error decoding JSON: %s", err)
				continue
			}

			headerBytes, err := json.Marshal(d.Headers)
			if err != nil {
				log.Fatalf("Error encoding headers to JSON: %s", err)
				continue
			}

			headersMap := make(map[string]interface{})
			err = json.Unmarshal(headerBytes, &headersMap)
			if err != nil {
				log.Fatalf("Error decoding headers JSON: %s", err)
				continue
			}

			consumer.Consume(bodyMap, headersMap)
		}
	}()
}
