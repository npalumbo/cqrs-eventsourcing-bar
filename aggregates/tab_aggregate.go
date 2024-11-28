package aggregates

import (
	"errors"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/domain"
	"golangsevillabar/events"
	"golangsevillabar/utils"
	"slices"

	"github.com/thoas/go-funk"
)

type tabAggregate struct {
	tabOpen           bool
	eventEmitter      events.EventEmitter[events.Event]
	outstandingDrinks []domain.OrderedItem
	servedItemsAmount float64
}

type TabAggregate interface {
	HandleCommand(c commands.Command) error
	ApplyEvent(e events.Event) error
}

func (t tabAggregate) HandleCommand(c commands.Command) error {
	switch command := c.(type) {
	case commands.OpenTab:
		return t.handleCommandOpenTab(command)
	case commands.PlaceOrder:
		return t.handleCommandPlaceOrder(command)
	case commands.MarkDrinksServed:
		return t.handleCommandMarkDrinksServed(command)
	case commands.CloseTab:
		return t.handleCommandCloseTab(command)
	default:
		return fmt.Errorf("unexpected commands.Command: %#v", c)
	}
}

func (t *tabAggregate) ApplyEvent(e events.Event) error {
	switch event := e.(type) {
	case events.TabOpened:
		return t.applyTabOpened(event)
	case events.DrinksOrdered:
		return t.applyDrinksOrdered(event)
	case events.DrinkServed:
		return t.applyDrinksServed(event)
	case events.TabClosed:
		return t.applyTabClosed(event)
	default:
		return fmt.Errorf("unexpected events.Event: %#v", e)
	}
}

func (t *tabAggregate) handleCommandOpenTab(c commands.OpenTab) error {
	t.eventEmitter.EmitEvent(events.TabOpened{ID: c.ID, TableNumber: c.TableNumber, Waiter: c.Waiter})
	return nil
}

func (t *tabAggregate) handleCommandPlaceOrder(c commands.PlaceOrder) error {
	if t.tabOpen {
		t.eventEmitter.EmitEvent(events.DrinksOrdered{ID: c.ID, Items: c.Items})
		return nil
	}
	return errors.New("tab is not opened")
}

func (t *tabAggregate) handleCommandMarkDrinksServed(c commands.MarkDrinksServed) error {
	menuItemsThatAreNotInOrderedItems := utils.FindMenuItemsThatAreNotInOrderedItems(t.outstandingDrinks, c.MenuNumbers)
	if len(menuItemsThatAreNotInOrderedItems) > 0 {
		return fmt.Errorf("cannot serve drinks that were not ordered: %v", menuItemsThatAreNotInOrderedItems)
	}

	t.eventEmitter.EmitEvent(events.DrinkServed{ID: c.ID, MenuNumbers: c.MenuNumbers})
	return nil
}

func (t *tabAggregate) handleCommandCloseTab(c commands.CloseTab) error {
	servedItemsAmount := t.servedItemsAmount
	if !t.tabOpen {
		return errors.New("cannot close a tab that is not open")
	}
	if len(t.outstandingDrinks) > 0 {
		return errors.New("cannot close a tab with unserved items")
	}
	if c.AmountPaid < servedItemsAmount {
		return fmt.Errorf("not enough to cover tab, total served cost is: %v, but paid: %v", servedItemsAmount, c.AmountPaid)
	}
	t.eventEmitter.EmitEvent(events.TabClosed{ID: c.ID, AmountPaid: c.AmountPaid, OrderAmount: servedItemsAmount, Tip: c.AmountPaid - servedItemsAmount})
	return nil

}

func (t *tabAggregate) applyTabOpened(e events.TabOpened) error {
	t.tabOpen = true
	return nil
}

func (t *tabAggregate) applyDrinksOrdered(e events.DrinksOrdered) error {
	t.outstandingDrinks = e.Items
	return nil
}

func (t *tabAggregate) applyDrinksServed(e events.DrinkServed) error {
	for _, menuNumber := range e.MenuNumbers {
		found := funk.Find(t.outstandingDrinks, func(item domain.OrderedItem) bool { return item.MenuItem == menuNumber })
		if found != nil {
			if itemFound, ok := found.(domain.OrderedItem); ok {
				t.outstandingDrinks = slices.DeleteFunc(t.outstandingDrinks, func(itemToDelete domain.OrderedItem) bool {
					return itemToDelete == itemFound
				})
				t.servedItemsAmount += itemFound.Price
			}

		}
	}

	return nil
}

func (t *tabAggregate) applyTabClosed(e events.TabClosed) error {
	t.tabOpen = false
	return nil
}

func CreateTabAggregate(eventEmitter events.EventEmitter[events.Event]) TabAggregate {
	return &tabAggregate{
		tabOpen:           false,
		eventEmitter:      eventEmitter,
		outstandingDrinks: []domain.OrderedItem{},
		servedItemsAmount: 0,
	}
}
