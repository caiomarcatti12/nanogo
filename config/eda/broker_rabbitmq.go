package eda

import (
	"github.com/caiomarcatti12/nanogo/v2/config/rabbitmq"
)

type RabbitMQBroker struct {
}

func NewEdaRabbitmq() *RabbitMQBroker {
	return &RabbitMQBroker{}
}

func (b *RabbitMQBroker) SendEvent(event Event) {
	rabbitmq.Publish(event.Channel, event.Key, event.Data)
}
