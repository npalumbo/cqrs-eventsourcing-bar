package aggregates_test

import (
	"golangsevillabar/aggregates"
	"golangsevillabar/commands"
	"golangsevillabar/domain"
	"golangsevillabar/events"
	mock_events "golangsevillabar/mocks/events"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTabAggregate_CanOpenTab(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	eventEmitter := mock_events.NewMockEventEmitter[events.Event](mockCtrl)
	tabAggregate := aggregates.CreateTabAggregate(eventEmitter)

	commandID, _ := ksuid.NewRandom()

	eventEmitter.EXPECT().EmitEvent(events.TabOpened{
		ID:          commandID,
		TableNumber: 0,
		Waiter:      "waiter_1",
	}).Times(1)

	err := tabAggregate.OpenTabHandler.Handle(commands.OpenTab{
		ID:          commandID,
		TableNumber: 0,
		Waiter:      "waiter_1",
	})

	assert.NoError(t, err)
}

func TestTabAggregate_CanNotOrderWithUnOpenedTab(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	eventEmitter := mock_events.NewMockEventEmitter[events.Event](mockCtrl)
	tabAggregate := aggregates.CreateTabAggregate(eventEmitter)

	commandID, _ := ksuid.NewRandom()

	eventEmitter.EXPECT().EmitEvent(events.DrinksOrdered{
		ID:    commandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	}).Times(0)

	err := tabAggregate.PlaceOrderHandler.Handle(commands.PlaceOrder{
		ID:    commandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})

	assert.Error(t, err)
	assert.Equal(t, "tab is not opened", err.Error())
}

func TestTabAggregate_CanOrderWhenTabIsOpen(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	eventEmitter := mock_events.NewMockEventEmitter[events.Event](mockCtrl)
	tabAggregate := aggregates.CreateTabAggregate(eventEmitter)

	tabOpenedEventID, _ := ksuid.NewRandom()
	placeOrderCommandID, _ := ksuid.NewRandom()

	tabAggregate.TabOpenedHandler.Apply(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})

	eventEmitter.EXPECT().EmitEvent(events.DrinksOrdered{
		ID:    placeOrderCommandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	}).Times(1)

	err := tabAggregate.PlaceOrderHandler.Handle(commands.PlaceOrder{
		ID:    placeOrderCommandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})

	assert.NoError(t, err)
}
