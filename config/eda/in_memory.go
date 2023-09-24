package eda

type InMemoryBroker struct {
}

func NewEdaInMemory() *InMemoryBroker {
	return &InMemoryBroker{}
}

func (b *InMemoryBroker) SendEvent(event Event) {
	// TODO
}
