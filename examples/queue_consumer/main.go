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

	exchange := queue.RabbitmqExchange{
		Name:    "wrk-log-trail:event:logged",
		Type:    queue.Direct,
		Durable: true,
	}

	queueCfg := queue.RabbitmqQueue{
		Name:        "wrk-log-trail:event:logged:consumer",
		RoutingKey:  "trail_event_logged",
		ConsumerTag: "trail_event_logged_consumer",
		Durable:     true,
	}

	_ = queueManager.Configure(exchange, queueCfg)

	consumer := NewEventConsumer(logger)
	err = queueManager.Consume(&queueCfg, consumer)
	if err != nil {
		panic(err)
	}

	nanogo.WaitSignalStop()
}
