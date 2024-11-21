package aggregates_test

import (
	"golangsevillabar/aggregates"
	"golangsevillabar/commands"
	"golangsevillabar/domain"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestTabAggregate_CanOpenTab(t *testing.T) {
	tabAggregate := aggregates.CreateTabAggregate()

	commandID, _ := ksuid.NewRandom()

	event, err := tabAggregate.OpenTabHandler.HandleCommand(commands.OpenTab{
		ID:          commandID,
		TableNumber: 0,
		Waiter:      "waiter_1",
	})

	assert.Nil(t, err)
	assert.Equal(t, commandID, event.ID)
	assert.Equal(t, 0, event.TableNumber)
	assert.Equal(t, "waiter_1", event.Waiter)
}

func TestTabAggregate_CanNotOrderWithUnOpenedTab(t *testing.T) {
	tabAggregate := aggregates.CreateTabAggregate()

	commandID, _ := ksuid.NewRandom()

	_, err := tabAggregate.PlaceOrderHandler.HandleCommand(commands.PlaceOrder{
		ID:    commandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})

	assert.Error(t, err)
	assert.Equal(t, "Tab is not opened", err.Error())
}
