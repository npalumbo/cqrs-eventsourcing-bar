package commands

import (
	"context"
	"fmt"
	"golangsevillabar/events"
	"reflect"
)

//go:generate mockery --name CommandDispatcher
type CommandDispatcher interface {
	DispatchCommand(ctx context.Context, command Command) error
}

type Dispatcher struct {
	eventStore       events.EventStore
	eventEmitter     events.EventEmitter
	aggregateFactory AggregateFactory
}

func CreateCommandDispatcher(eventStore events.EventStore, eventEmitter events.EventEmitter, aggregateFactory AggregateFactory) *Dispatcher {
	return &Dispatcher{eventStore: eventStore, eventEmitter: eventEmitter, aggregateFactory: aggregateFactory}
}

func (d *Dispatcher) DispatchCommand(ctx context.Context, command Command) error {
	aggregate := d.aggregateFactory.CreateAggregate()

	eventsLoaded, err := d.eventStore.LoadEvents(ctx, command.GetID())

	if err != nil {
		return fmt.Errorf("error loading events for aggregate: %s, reason: %w", command.GetID().String(), err)
	}

	for i, event := range eventsLoaded {
		err = aggregate.ApplyEvent(event)
		if err != nil {
			return fmt.Errorf("error applying past event [%s-#%d] for aggregate: %s, reason: %w", events.GetEventTypeAsString(event), i, command.GetID().String(), err)
		}
	}

	newEvents, err := aggregate.HandleCommand(command)

	if err != nil {
		return fmt.Errorf("error handling command [%s] for aggregate: %s, reason: %w", reflect.TypeOf(command).Name(), command.GetID().String(), err)
	}

	err = d.eventStore.SaveEvents(ctx, command.GetID(), len(eventsLoaded), newEvents)

	if err != nil {
		return fmt.Errorf("error when saving events for aggregate: %s, reason: %w", command.GetID(), err)
	}

	for _, event := range newEvents {
		err = d.eventEmitter.EmitEvent(event)
		if err != nil {
			return fmt.Errorf("error when emitting event [%s] for aggregate: %s, reason: %w", events.GetEventTypeAsString(event), command.GetID(), err)
		}
	}

	return nil
}
