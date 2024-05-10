package event

type IEventDispatcher interface {
	RegisterConsumer(eventKey string, handler EventHandler)
	Dispatch(event Event)
}
