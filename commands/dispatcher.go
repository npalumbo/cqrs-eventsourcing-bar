package commands

import (
	"fmt"
	"golangsevillabar/events"
	"reflect"
)

type Dispatcher struct {
	eventStore   events.EventStore
	eventEmitter events.EventEmitter
}

func CreateCommandDispatcher(eventStore events.EventStore, eventEmitter events.EventEmitter) *Dispatcher {
	return &Dispatcher{eventStore: eventStore, eventEmitter: eventEmitter}
}

func (d *Dispatcher) DispatchCommand(command Command) error {
	aggregate := CreateTabAggregate()

	events, err := d.eventStore.LoadEvents(command.GetID())

	if err != nil {
		return fmt.Errorf("error loading events for aggregate: %d, %w", command.GetID(), err)
	}

	for i, event := range events {
		err = aggregate.ApplyEvent(event)
		if err != nil {
			return fmt.Errorf("error applying past event [%s-#%d] for aggregate: %d, %w", reflect.TypeOf(event).Name(), i, command.GetID(), err)
		}
	}

	newEvents, err := aggregate.HandleCommand(command)

	if err != nil {
		return fmt.Errorf("error handling command [%s] for aggregate: %d, %w", reflect.TypeOf(command).Name(), command.GetID(), err)
	}

	err = d.eventStore.SaveEvents(command.GetID(), len(events), newEvents)

	if err != nil {
		return fmt.Errorf("error when saving events for aggregate: %d, %w", command.GetID(), err)
	}

	for _, event := range newEvents {
		d.eventEmitter.EmitEvent(event)
	}

	return nil
}
