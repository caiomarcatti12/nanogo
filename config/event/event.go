package event

type Event struct {
	Channel string
	Key     string
	Data    interface{}
}
