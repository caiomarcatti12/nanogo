package queue

import "github.com/caiomarcatti12/nanogo/v2/config/rabbitmq"

func Exchange() rabbitmq.Exchange {
	return rabbitmq.Exchange{
		Name:    "tax-invoice",
		Durable: true,
		Type:    "topic",
		AutoDel: false,
		NoWait:  false,
	}
}

func QueueIntegrationCityHallCancel() rabbitmq.Queue {
	return rabbitmq.Queue{
		Name:       "cancel",
		Key:        "cancel",
		Durable:    true,
		AutoDel:    false,
		Exclusive:  false,
		NoWait:     false,
		Parameters: nil,
	}
}

func QueueIntegrationCityHallRegister() rabbitmq.Queue {
	return rabbitmq.Queue{
		Name:       "register",
		Key:        "register",
		Durable:    true,
		AutoDel:    false,
		Exclusive:  false,
		NoWait:     false,
		Parameters: nil,
	}
}

func QueueIntegrationCityHallWaitingRegister() rabbitmq.Queue {
	return rabbitmq.Queue{
		Name:       "waiting",
		Key:        "waiting",
		Durable:    true,
		AutoDel:    false,
		Exclusive:  false,
		NoWait:     false,
		Parameters: nil,
	}
}
