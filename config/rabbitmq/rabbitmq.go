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
	"fmt"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/config/env"
	logger "github.com/caiomarcatti12/nanogo/v2/config/log"

	"github.com/streadway/amqp"
)

type Connection struct {
	*amqp.Connection
	Channel *amqp.Channel
}

var instance *Connection
var once sync.Once

func NewInstanceRabbitmq() *Connection {
	once.Do(func() {
		rabbitmqUser := env.GetEnv("RABBITMQ_USER")
		rabbitmqPassword := env.GetEnv("RABBITMQ_PASSWORD")
		rabbitmqHost := env.GetEnv("RABBITMQ_HOST")
		rabbitmqPort := env.GetEnv("RABBITMQ_PORT")
		rabbitmqVhost := env.GetEnv("RABBITMQ_VHOST")

		url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort, rabbitmqVhost)

		conn, err := amqp.Dial(url)

		if err != nil {
			logger.Fatal("Failed to connect to RabbitMQ: %s", err)
		}

		ch, err := conn.Channel()

		if err != nil {
			logger.Fatal("Failed to open a channel: %s", err)
		}

		instance = &Connection{conn, ch}
	})

	return instance
}
