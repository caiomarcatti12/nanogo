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

type IDemoConsumer interface {
	Handler(msg DemoMessage, headers map[string]interface{}) error
}

type DemoConsumer struct {
	logger log.ILog
}

func NewDemoConsumer(logger log.ILog) IDemoConsumer {
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

	// Configuração da fila
	queueCfg := queue.NatsQueue{
		Name:       "demo.subject",
		QueueGroup: "demo-group",
	}

	// Novo método AddConsumer - simplifica todo o processo!
	consumer := queue.QueueConsumer{
		Queue:   &queueCfg,
		Handler: NewDemoConsumer,
	}

	// Registra o consumer, configura a fila e inicia o consumo em uma única chamada
	if err := queueManager.AddConsumer(consumer); err != nil {
		panic(err)
	}

	// Publica uma mensagem de teste
	msg := DemoMessage{ID: "1", Text: "hello from nanogo with new AddConsumer!"}
	if err := queueManager.Publish(queueCfg.Name, "", msg); err != nil {
		logger.Error("publish failed", "err", err)
	}

	nanogo.WaitSignalStop()
}
