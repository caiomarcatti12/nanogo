package main

import (
	"github.com/caiomarcatti12/nanogo/config"
	"github.com/caiomarcatti12/nanogo/config/env"
	"github.com/caiomarcatti12/nanogo/config/rabbitmq"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	//Cria um consumidor da fila MyConsumer
	rabbitmq.Consume(exchange(), queue(), &MyConsumer{})

	config.WaitSignalStop()
}

func exchange() rabbitmq.Exchange {
	return rabbitmq.Exchange{
		Name:    "teste-exchange-go",
		Durable: true,
		Type:    "direct",
		AutoDel: false,
		NoWait:  false,
	}
}

func queue() rabbitmq.Queue {
	return rabbitmq.Queue{
		Name:       "teste-queue",
		Durable:    true,
		AutoDel:    false,
		Exclusive:  false,
		NoWait:     false,
		Parameters: nil,
	}
}
