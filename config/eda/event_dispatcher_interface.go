package eda

type IEventDispatcher interface {
	SendEvent(event Event)
}
