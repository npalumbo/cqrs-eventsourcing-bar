package events

//go:generate mockery --name EventListener
type EventListener interface {
	HandleEvent(e Event) error
}
