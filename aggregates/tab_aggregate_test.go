package aggregates_test

import (
	"golangsevillabar/aggregates"
	"golangsevillabar/commands"
	"golangsevillabar/domain"

	// "golangsevillabar/domain"
	"golangsevillabar/events"
	"golangsevillabar/events/mocks"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TabAggregateTestSuite struct {
	suite.Suite
	tabAggregate *aggregates.TabAggregate
	eventEmitter *mocks.EventEmitter[events.Event]
}

func (suite *TabAggregateTestSuite) SetupTest() {
	suite.eventEmitter = mocks.NewEventEmitter[events.Event](suite.T())
	tabAggregate := aggregates.CreateTabAggregate(suite.eventEmitter)
	suite.tabAggregate = &tabAggregate
	suite.eventEmitter.On("EmitEvent", mock.Anything).Maybe()
}

func (suite *TabAggregateTestSuite) TestCanOpenTab() {

	commandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given

	// When
	err := suite.tabAggregate.HandleCommand(commands.OpenTab{
		ID:          commandID,
		TableNumber: 0,
		Waiter:      "waiter_1",
	})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.TabOpened{
		ID:          commandID,
		TableNumber: 0,
		Waiter:      "waiter_1",
	})
}

func (suite *TabAggregateTestSuite) TestCanNotOrderWithUnOpenedTab() {
	commandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given

	// When
	err := suite.tabAggregate.HandleCommand(commands.PlaceOrder{
		ID:    commandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "tab is not opened", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func (suite *TabAggregateTestSuite) TestCanOrderWhenTabIsOpen() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	placeOrderCommandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	// Given(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})

	// When
	err := suite.tabAggregate.HandleCommand(commands.PlaceOrder{
		ID:    placeOrderCommandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.DrinksOrdered{
		ID:    placeOrderCommandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "cruzcampo", Price: 1.5}},
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TabAggregateTestSuite))
}
