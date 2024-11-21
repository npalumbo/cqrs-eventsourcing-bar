package aggregates

import (
	"errors"
	"golangsevillabar/commands"
	"golangsevillabar/events"
)

type CommandHandler[T commands.Command, E any] interface {
	HandleCommand(command T) (E, error)
}

type openTabHandler struct {
}

type placeOrderHandler struct {
	tabOpen *bool
}

func (o openTabHandler) HandleCommand(command commands.OpenTab) (*events.TabOpened, error) {
	return &events.TabOpened{
		ID:          command.ID,
		TableNumber: command.TableNumber,
		Waiter:      command.Waiter,
	}, nil
}

func (p placeOrderHandler) HandleCommand(command commands.PlaceOrder) (*events.DrinksOrdered, error) {
	if *p.tabOpen {
		return &events.DrinksOrdered{
			ID:    command.ID,
			Items: command.Items,
		}, nil
	}
	return nil, errors.New("Tab is not opened")

}

type TabAggregate struct {
	tabOpen           *bool
	OpenTabHandler    CommandHandler[commands.OpenTab, *events.TabOpened]
	PlaceOrderHandler CommandHandler[commands.PlaceOrder, *events.DrinksOrdered]
}

func CreateTabAggregate() TabAggregate {
	tabOpen := false
	return TabAggregate{
		tabOpen:           &tabOpen,
		OpenTabHandler:    openTabHandler{},
		PlaceOrderHandler: placeOrderHandler{&tabOpen},
	}
}
