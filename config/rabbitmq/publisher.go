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

	logger "github.com/caiomarcatti12/nanogo/v2/config/log"
	"github.com/streadway/amqp"
)

func Publish(exchangeName string, routingKey string, body map[string]interface{}) {
	connection := NewInstanceRabbitmq()

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		logger.Fatalf("Houve uma falha ao converter a struct para json: %s", err)
	}

	errPublish := connection.Channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)

	if errPublish != nil {
		logger.Fatalf("Failed to publish a message: %s", errPublish)
	}
}
