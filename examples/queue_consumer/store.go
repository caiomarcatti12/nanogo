package main

import (
	"encoding/json"
	"fmt"
	"github.com/caiomarcatti12/nanogo/pkg/log"
)

type EventConsumer struct {
	logger log.ILog
}

func NewEventConsumer(logger log.ILog) *EventConsumer {
	return &EventConsumer{
		logger: logger,
	}
}

func (c *EventConsumer) Handler(event Event, headers map[string]interface{}) error {
	if event.ID == "" {
		return fmt.Errorf("evento com ID vazio não pode ser processado")
	}

	// Realiza o marshal do evento usando json.Marshal
	jsonData, err := json.Marshal(event)
	if err != nil {
		// Log de erro caso o marshal falhe
		c.logger.Error("Erro ao marshallizar o evento", "error", err)
		return err
	}

	// Log do JSON gerado
	c.logger.Info("Evento processado", "data", string(jsonData))
	return nil
}
