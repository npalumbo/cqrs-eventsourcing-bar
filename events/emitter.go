package events

import "fmt"

//go:generate mockery --name EventEmitter
type EventEmitter interface {
	EmitEvent(event Event) error
}

type PrintlnEventEmitter struct {
}

func (e PrintlnEventEmitter) EmitEvent(event Event) error {
	fmt.Printf("emitted: %v\n", event)
	return nil
}
