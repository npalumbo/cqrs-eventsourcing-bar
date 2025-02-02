package commands

import (
	"cqrseventsourcingbar/shared"

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
	TableNumber int
	Waiter      string
}

type PlaceOrder struct {
	BaseCommand
	Items []shared.MenuItem
}

type MarkDrinksServed struct {
	BaseCommand
	MenuNumbers []int
}

type CloseTab struct {
	BaseCommand
	AmountPaid float64
}
