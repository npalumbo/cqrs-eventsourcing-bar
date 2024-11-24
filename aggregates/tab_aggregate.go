package aggregates

import (
	"errors"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/domain"
	"golangsevillabar/events"
	"golangsevillabar/utils"
)

type CommandHandler[T commands.Command] interface {
	handle(c T) error
}

type EventApplier[E events.Event] interface {
	apply(e E) error
}

// Command Handlers
type openTabHandler struct {
	eventEmitter events.EventEmitter[events.Event]
}
type placeOrderHandler struct {
	tabAggregate *TabAggregate
	eventEmitter events.EventEmitter[events.Event]
}
type markDrinksServedHandler struct {
	tabAggregate *TabAggregate
	eventEmitter events.EventEmitter[events.Event]
}

// Event Appliers
type tabOpenedApplier struct {
	tabAggregate *TabAggregate
}
type drinksOrderedApplier struct {
	tabAggregate *TabAggregate
}
type drinksServedApplier struct {
	tabAggregate *TabAggregate
}

type TabAggregate struct {
	tabOpen           bool
	outstandingDrinks []domain.OrderedItem
	// Command Handlers
	openTabHandler          CommandHandler[commands.OpenTab]
	placeOrderHandler       CommandHandler[commands.PlaceOrder]
	markDrinksServedHandler CommandHandler[commands.MarkDrinksServed]
	// Event Appliers
	tabOpenedApplier     EventApplier[events.TabOpened]
	drinksOrderedApplier EventApplier[events.DrinksOrdered]
	drinksServedApplier  EventApplier[events.DrinkServed]
}

func (t TabAggregate) HandleCommand(c commands.Command) error {
	switch command := c.(type) {
	case commands.OpenTab:
		return t.openTabHandler.handle(command)
	case commands.PlaceOrder:
		return t.placeOrderHandler.handle(command)
	case commands.MarkDrinksServed:
		return t.markDrinksServedHandler.handle(command)
	default:
		return fmt.Errorf("unexpected commands.Command: %#v", c)
	}
}

func (t TabAggregate) ApplyEvent(e events.Event) error {
	switch event := e.(type) {
	case events.TabOpened:
		return t.tabOpenedApplier.apply(event)
	case events.DrinksOrdered:
		return t.drinksOrderedApplier.apply(event)
	case events.DrinkServed:
		return t.drinksServedApplier.apply(event)
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

func (m markDrinksServedHandler) handle(c commands.MarkDrinksServed) error {
	menuItemsThatAreNotInOrderedItems := utils.FindMenuItemsThatAreNotInOrderedItems(m.tabAggregate.outstandingDrinks, c.MenuNumbers)
	if len(menuItemsThatAreNotInOrderedItems) > 0 {
		return fmt.Errorf("cannot serve drinks that were not ordered: %v", menuItemsThatAreNotInOrderedItems)
	}

	m.eventEmitter.EmitEvent(events.DrinkServed{ID: c.ID, MenuNumbers: c.MenuNumbers})
	return nil
}

func (t tabOpenedApplier) apply(e events.TabOpened) error {
	t.tabAggregate.tabOpen = true
	return nil
}

func (d drinksOrderedApplier) apply(e events.DrinksOrdered) error {
	d.tabAggregate.outstandingDrinks = e.Items
	return nil
}

func (d drinksServedApplier) apply(e events.DrinkServed) error {
	d.tabAggregate.outstandingDrinks = utils.RemoveOrderedItemsThatAppearInMarkedServedItems(d.tabAggregate.outstandingDrinks, e.MenuNumbers)
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
		markDrinksServedHandler: markDrinksServedHandler{
			tabAggregate: &tabAggregate,
			eventEmitter: eventEmitter,
		},
		tabOpenedApplier: tabOpenedApplier{
			tabAggregate: &tabAggregate,
		},
		drinksOrderedApplier: drinksOrderedApplier{
			tabAggregate: &tabAggregate,
		},
		drinksServedApplier: drinksServedApplier{
			tabAggregate: &tabAggregate,
		},
	}
}
