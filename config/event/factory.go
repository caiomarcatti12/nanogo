package event

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

func FactoryEventDispatcher(env env.IEnv, log log.ILog) IEventDispatcher {
	eventDispatcher := env.GetEnv("EVENT_DISPATCHER", "IN_MEMORY")

	switch eventDispatcher {
	case "IN_MEMORY":
		return NewInMemoryBroker(log)
	default:
		return NewInMemoryBroker(log)
	}
}
