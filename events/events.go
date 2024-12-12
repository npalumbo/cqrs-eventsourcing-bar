package events

import (
	"golangsevillabar/shared"

	"github.com/segmentio/ksuid"
)

type Event interface{}

type TabOpened struct {
	ID          ksuid.KSUID
	TableNumber int
	Waiter      string
}

type DrinksOrdered struct {
	ID    ksuid.KSUID
	Items []shared.OrderedItem
}

type DrinkServed struct {
	ID          ksuid.KSUID
	MenuNumbers []int
}

type TabClosed struct {
	ID          ksuid.KSUID
	AmountPaid  float64
	OrderAmount float64
	Tip         float64
}
