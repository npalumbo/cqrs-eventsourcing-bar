package commands

import (
	"golangsevillabar/domain"

	"github.com/segmentio/ksuid"
)

type Command interface {
	GetID() ksuid.KSUID
}

type BaseCommand struct {
	ID ksuid.KSUID
}

func (command BaseCommand) GetID() ksuid.KSUID {
	return command.ID
}

type OpenTab struct {
	BaseCommand
	ID          ksuid.KSUID
	TableNumber int
	Waiter      string
}

type PlaceOrder struct {
	BaseCommand
	ID    ksuid.KSUID
	Items []domain.OrderedItem
}

type MarkDrinksServed struct {
	BaseCommand
	ID          ksuid.KSUID
	MenuNumbers []int
}

type CloseTab struct {
	BaseCommand
	ID         ksuid.KSUID
	AmountPaid float64
}
