package aggregates

import (
	"errors"
	"golangsevillabar/commands"
	"golangsevillabar/events"
)

type CommandHandler[T commands.Command] interface {
	Handle(c T) error
}

type EventHandler[E events.Event] interface {
	Apply(e E)
}

// Commands
type OpenTabHandler CommandHandler[commands.OpenTab]
type PlaceOrderHandler CommandHandler[commands.PlaceOrder]

// Events
type TabOpenedHandler EventHandler[events.TabOpened]

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

func (t openTabHandler) Handle(c commands.OpenTab) error {
	t.eventEmitter.EmitEvent(events.TabOpened{ID: c.ID, TableNumber: c.TableNumber, Waiter: c.Waiter})
	return nil
}

func (p placeOrderHandler) Handle(c commands.PlaceOrder) error {
	if p.tabAggregate.tabOpen {
		p.eventEmitter.EmitEvent(events.DrinksOrdered{ID: c.ID, Items: c.Items})
		return nil
	}
	return errors.New("tab is not opened")
}

func (t tabOpenedApplier) Apply(e events.TabOpened) {
	t.tabAggregate.tabOpen = true
}

type TabAggregate struct {
	tabOpen           bool
	OpenTabHandler    OpenTabHandler
	PlaceOrderHandler PlaceOrderHandler
	TabOpenedHandler  TabOpenedHandler
}

func CreateTabAggregate(eventEmitter events.EventEmitter[events.Event]) (out TabAggregate) {
	return TabAggregate{
		tabOpen: false,
		OpenTabHandler: openTabHandler{
			eventEmitter: eventEmitter,
		},
		PlaceOrderHandler: placeOrderHandler{
			tabAggregate: &out,
			eventEmitter: eventEmitter,
		},
		TabOpenedHandler: tabOpenedApplier{
			tabAggregate: &out,
		},
	}
}
