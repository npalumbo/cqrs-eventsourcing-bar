package events

type EventListener interface {
	HandleEvent(e Event) error
}
