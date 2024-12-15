package commands

import (
	"fmt"
	"golangsevillabar/events"
	"reflect"
)

type Dispatcher struct {
	eventStore       events.EventStore
	eventEmitter     events.EventEmitter
	aggregateFactory AggregateFactory
}

func CreateCommandDispatcher(eventStore events.EventStore, eventEmitter events.EventEmitter, aggregateFactory AggregateFactory) *Dispatcher {
	return &Dispatcher{eventStore: eventStore, eventEmitter: eventEmitter, aggregateFactory: aggregateFactory}
}

func (d *Dispatcher) DispatchCommand(command Command) error {
	aggregate := d.aggregateFactory.CreateAggregate()

	events, err := d.eventStore.LoadEvents(command.GetID())

	if err != nil {
		return fmt.Errorf("error loading events for aggregate: %s, reason: %w", command.GetID().String(), err)
	}

	for i, event := range events {
		err = aggregate.ApplyEvent(event)
		if err != nil {
			return fmt.Errorf("error applying past event [%s-#%d] for aggregate: %s, reason: %w", reflect.TypeOf(event).Name(), i, command.GetID().String(), err)
		}
	}

	newEvents, err := aggregate.HandleCommand(command)

	if err != nil {
		return fmt.Errorf("error handling command [%s] for aggregate: %s, reason: %w", reflect.TypeOf(command).Name(), command.GetID().String(), err)
	}

	err = d.eventStore.SaveEvents(command.GetID(), len(events), newEvents)

	if err != nil {
		return fmt.Errorf("error when saving events for aggregate: %s, reason: %w", command.GetID(), err)
	}

	for _, event := range newEvents {
		err = d.eventEmitter.EmitEvent(event)
		if err != nil {
			return fmt.Errorf("error when emitting event [%s] for aggregate: %s, reason: %w", reflect.TypeOf(event), command.GetID(), err)
		}
	}

	return nil
}
