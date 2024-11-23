package events

import "fmt"

//go:generate mockgen -destination=../mocks/events/mock_emitter.go -source=./emitter.go
type EventEmitter[E Event] interface {
	EmitEvent(event E)
}

type PrintlnEventEmitter[E Event] struct {
}

func (e PrintlnEventEmitter[E]) EmitEvent(event E) {
	fmt.Printf("emitted: %v\n", event)
}
