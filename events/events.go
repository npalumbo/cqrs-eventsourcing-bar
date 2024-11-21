package events

import (
	"golangsevillabar/domain"

	"github.com/segmentio/ksuid"
)

type Event interface{}

type TabOpened struct {
	Event
	ID          ksuid.KSUID
	TableNumber int
	Waiter      string
}

type DrinksOrdered struct {
	Event
	ID    ksuid.KSUID
	Items []domain.OrderedItem
}
