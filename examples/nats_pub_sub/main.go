package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/caiomarcatti12/nanogo/pkg/nanogo"
	"github.com/caiomarcatti12/nanogo/pkg/queue"
)

type DemoMessage struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type DemoConsumer struct {
	logger log.ILog
}

func NewDemoConsumer(logger log.ILog) *DemoConsumer {
	return &DemoConsumer{logger: logger}
}

func (c *DemoConsumer) Handler(msg DemoMessage, headers map[string]interface{}) error {
	c.logger.Info("received message", "id", msg.ID, "text", msg.Text)
	return nil
}

func main() {
	nanogo.Bootstrap()

	queueManager, err := di.Get[queue.IQueue]()
	if err != nil {
		panic(err)
	}

	logger, _ := di.Get[log.ILog]()

	queueCfg := queue.NatsQueue{
		Name:       "demo.subject",
		QueueGroup: "demo-group",
	}

	_ = queueManager.Configure(queueCfg)

	consumer := NewDemoConsumer(logger)

	go func() {
		if err := queueManager.Consume(&queueCfg, consumer); err != nil {
			logger.Error("error consuming", "err", err)
		}
	}()

	msg := DemoMessage{ID: "1", Text: "hello from nanogo"}
	if err := queueManager.Publish(queueCfg.Name, "", msg); err != nil {
		logger.Error("publish failed", "err", err)
	}

	nanogo.WaitSignalStop()
}
