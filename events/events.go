package events

import (
	"golangsevillabar/shared"

	"github.com/segmentio/ksuid"
)

type Event interface{}

type BaseEvent struct {
	ID ksuid.KSUID
}

func (event BaseEvent) GetID() ksuid.KSUID {
	return event.ID
}

type TabOpened struct {
	BaseEvent
	TableNumber int
	Waiter      string
}

type DrinksOrdered struct {
	BaseEvent
	Items []shared.OrderedItem
}

type DrinkServed struct {
	BaseEvent
	MenuNumbers []int
}

type TabClosed struct {
	BaseEvent
	AmountPaid  float64
	OrderAmount float64
	Tip         float64
}
