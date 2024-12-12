package commands

import (
	"errors"
	"fmt"
	"golangsevillabar/events"
	"golangsevillabar/shared"
	"slices"

	"github.com/thoas/go-funk"
)

type tabAggregate struct {
	tabOpen           bool
	outstandingDrinks []shared.OrderedItem
	servedItemsAmount float64
}

type TabAggregate interface {
	HandleCommand(c Command) ([]events.Event, error)
	ApplyEvent(e events.Event) error
}

func (t tabAggregate) HandleCommand(c Command) ([]events.Event, error) {
	switch command := c.(type) {
	case OpenTab:
		return t.handleCommandOpenTab(command)
	case PlaceOrder:
		return t.handleCommandPlaceOrder(command)
	case MarkDrinksServed:
		return t.handleCommandMarkDrinksServed(command)
	case CloseTab:
		return t.handleCommandCloseTab(command)
	default:
		return nil, fmt.Errorf("unexpected Command: %#v", c)
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

func (t *tabAggregate) handleCommandOpenTab(c OpenTab) ([]events.Event, error) {
	return []events.Event{events.TabOpened{ID: c.ID, TableNumber: c.TableNumber, Waiter: c.Waiter}}, nil
}

func (t *tabAggregate) handleCommandPlaceOrder(c PlaceOrder) ([]events.Event, error) {
	if t.tabOpen {
		return []events.Event{events.DrinksOrdered{ID: c.ID, Items: c.Items}}, nil
	}
	return nil, errors.New("tab is not opened")
}

func (t *tabAggregate) handleCommandMarkDrinksServed(c MarkDrinksServed) ([]events.Event, error) {
	menuItemsThatAreNotInOrderedItems := FindMenuItemsThatAreNotInOrderedItems(t.outstandingDrinks, c.MenuNumbers)
	if len(menuItemsThatAreNotInOrderedItems) > 0 {
		return nil, fmt.Errorf("cannot serve drinks that were not ordered: %v", menuItemsThatAreNotInOrderedItems)
	}

	return []events.Event{events.DrinkServed{ID: c.ID, MenuNumbers: c.MenuNumbers}}, nil
}

func (t *tabAggregate) handleCommandCloseTab(c CloseTab) ([]events.Event, error) {
	servedItemsAmount := t.servedItemsAmount
	if !t.tabOpen {
		return nil, errors.New("cannot close a tab that is not open")
	}
	if len(t.outstandingDrinks) > 0 {
		return nil, errors.New("cannot close a tab with unserved items")
	}
	if c.AmountPaid < servedItemsAmount {
		return nil, fmt.Errorf("not enough to cover tab, total served cost is: %v, but paid: %v", servedItemsAmount, c.AmountPaid)
	}
	return []events.Event{events.TabClosed{ID: c.ID, AmountPaid: c.AmountPaid, OrderAmount: servedItemsAmount, Tip: c.AmountPaid - servedItemsAmount}},
		nil

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
		found := funk.Find(t.outstandingDrinks, func(item shared.OrderedItem) bool { return item.MenuItem == menuNumber })
		if found != nil {
			if itemFound, ok := found.(shared.OrderedItem); ok {
				t.outstandingDrinks = slices.DeleteFunc(t.outstandingDrinks, func(itemToDelete shared.OrderedItem) bool {
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

func CreateTabAggregate() TabAggregate {
	return &tabAggregate{
		tabOpen:           false,
		outstandingDrinks: []shared.OrderedItem{},
		servedItemsAmount: 0,
	}
}
