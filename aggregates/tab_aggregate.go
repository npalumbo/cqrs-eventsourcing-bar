package aggregates

import (
	"errors"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/events"
)

type CommandHandler[T commands.Command] interface {
	handle(c T) error
}

type EventApplier[E events.Event] interface {
	apply(e E) error
}

type openTabHandler struct {
	eventEmitter events.EventEmitter[events.Event]
}
type placeOrderHandler struct {
	tabAggregate *TabAggregate
	eventEmitter events.EventEmitter[events.Event]
}
type tabOpenedApplier struct {
	tabAggregate *TabAggregate
}

type TabAggregate struct {
	tabOpen           bool
	openTabHandler    CommandHandler[commands.OpenTab]
	placeOrderHandler CommandHandler[commands.PlaceOrder]
	TabOpenedHandler  EventApplier[events.TabOpened]
}

func (t TabAggregate) HandleCommand(c commands.Command) error {
	switch command := c.(type) {
	case commands.OpenTab:
		return t.openTabHandler.handle(command)
	case commands.PlaceOrder:
		return t.placeOrderHandler.handle(command)
	default:
		return fmt.Errorf("unexpected commands.Command: %#v", c)
	}
}

func (t TabAggregate) ApplyEvent(e events.Event) error {
	switch event := e.(type) {
	case events.TabOpened:
		return t.TabOpenedHandler.apply(event)
	default:
		return fmt.Errorf("unexpected events.Event: %#v", e)
	}
}

func (t openTabHandler) handle(c commands.OpenTab) error {
	t.eventEmitter.EmitEvent(events.TabOpened{ID: c.ID, TableNumber: c.TableNumber, Waiter: c.Waiter})
	return nil
}

func (p placeOrderHandler) handle(c commands.PlaceOrder) error {
	if p.tabAggregate.tabOpen {
		p.eventEmitter.EmitEvent(events.DrinksOrdered{ID: c.ID, Items: c.Items})
		return nil
	}
	return errors.New("tab is not opened")
}

func (t tabOpenedApplier) apply(e events.TabOpened) error {
	t.tabAggregate.tabOpen = true
	return nil
}

func CreateTabAggregate(eventEmitter events.EventEmitter[events.Event]) (tabAggregate TabAggregate) {
	return TabAggregate{
		tabOpen: false,
		openTabHandler: openTabHandler{
			eventEmitter: eventEmitter,
		},
		placeOrderHandler: placeOrderHandler{
			tabAggregate: &tabAggregate,
			eventEmitter: eventEmitter,
		},
		TabOpenedHandler: tabOpenedApplier{
			tabAggregate: &tabAggregate,
		},
	}
}
