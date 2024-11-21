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
