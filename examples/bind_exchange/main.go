package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/nanogo"
	"github.com/caiomarcatti12/nanogo/pkg/queue"
)

func main() {
	nanogo.Bootstrap()

	queueManager, err := di.Get[queue.IQueue]()
	if err != nil {
		panic(err)
	}

	logger, _ := di.Get[log.ILog]()

	// Exchange principal
	mainExchange := queue.RabbitmqExchange{
		Name:    "wrk-log-trail:event:main",
		Type:    queue.Topic,
		Durable: true,
	}

	// Exchange secundária que será vinculada à principal
	secondaryExchange := queue.RabbitmqExchange{
		Name:    "wrk-log-trail:event:secondary",
		Type:    queue.Topic,
		Durable: true,
	}

	// Configuração da queue que receberá as mensagens da exchange secundária
	queueCfg := queue.RabbitmqQueue{
		Name:        "wrk-log-trail:event:secondary:consumer",
		RoutingKey:  "trail.event.*",
		ConsumerTag: "trail_event_secondary_consumer",
		Durable:     true,
	}

	// Configuração dos bindings
	bindingExchange := queue.BindingConfig{
		Source:      mainExchange.Name,
		Destination: secondaryExchange.Name,
		RoutingKey:  "trail.event.*",
		NoWait:      false,
		Args:        nil,
	}

	bindingQueue := queue.BindingConfig{
		Source:      secondaryExchange.Name,
		Destination: queueCfg.Name,
		RoutingKey:  "trail.event.*",
		NoWait:      false,
		Args:        nil,
	}

	// Configuração completa das exchanges, queue e bindings
	if err := queueManager.Configure(
		mainExchange,
		secondaryExchange,
		queueCfg,
		bindingExchange,
		bindingQueue,
	); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("RabbitMQ successfully configured and bindings established.")
}
