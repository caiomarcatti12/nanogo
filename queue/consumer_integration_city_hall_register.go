package queue

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

type Hello struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ConsumerIntegrationCityHallRegister struct{}

func (c *ConsumerIntegrationCityHallRegister) Consume(body Hello, headers map[string]interface{}) {
	log.Printf("Headers: %v", headers)
	fmt.Printf("Body: %s", body)
}
