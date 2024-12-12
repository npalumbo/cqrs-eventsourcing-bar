package events

import "fmt"

//go:generate mockery --name EventEmitter
type EventEmitter interface {
	EmitEvent(event Event)
}

type PrintlnEventEmitter struct {
}

func (e PrintlnEventEmitter) EmitEvent(event Event) {
	fmt.Printf("emitted: %v\n", event)
}
