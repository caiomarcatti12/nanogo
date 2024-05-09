package eda

import "github.com/caiomarcatti12/nanogo/v2/config/env"

func FactoryEventDispatcher(env env.IEnv) IEventDispatcher {
	eventDispatcher := env.GetEnv("EVENT_DISPATCHER", "IN_MEMORY")

	switch eventDispatcher {
	case "IN_MEMORY":
		return NewInMemoryBroker()
	default:
		return NewInMemoryBroker()
	}
}
