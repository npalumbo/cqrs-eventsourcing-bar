package commands

import (
	"golangsevillabar/domain"

	"github.com/segmentio/ksuid"
)

type Command interface{}

type OpenTab struct {
	ID          ksuid.KSUID
	TableNumber int
	Waiter      string
}

type PlaceOrder struct {
	ID    ksuid.KSUID
	Items []domain.OrderedItem
}

type MarkDrinksServed struct {
	ID          ksuid.KSUID
	MenuNumbers []int
}

type CloseTab struct {
	ID         ksuid.KSUID
	AmountPaid float64
}
