package events

import "fmt"

//go:generate mockery --name EventEmitter
type EventEmitter[E Event] interface {
	EmitEvent(event E)
}

type PrintlnEventEmitter[E Event] struct {
}

func (e PrintlnEventEmitter[E]) EmitEvent(event E) {
	fmt.Printf("emitted: %v\n", event)
}
