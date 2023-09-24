package eda

type EventDispatcherInterface interface {
	SendEvent(event Event)
}
