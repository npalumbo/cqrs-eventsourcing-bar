package aggregates_test

import (
	"golangsevillabar/aggregates"
	"golangsevillabar/commands"
	"golangsevillabar/domain"
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
	tabAggregate aggregates.TabAggregate
	eventEmitter *mocks.EventEmitter[events.Event]
}

func (suite *TabAggregateTestSuite) SetupTest() {
	suite.eventEmitter = mocks.NewEventEmitter[events.Event](suite.T())
	suite.tabAggregate = aggregates.CreateTabAggregate(suite.eventEmitter)
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
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "beer", Price: 1.5}},
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
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "beer", Price: 1.5}},
	})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.DrinksOrdered{
		ID:    placeOrderCommandID,
		Items: []domain.OrderedItem{{MenuItem: 11, Description: "beer", Price: 1.5}},
	})
}

func (suite *TabAggregateTestSuite) TestOrderedDrinksCanBeServed() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
		{MenuItem: 12, Description: "water", Price: 1.0},
	}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		ID:          markDrinksServedID,
		MenuNumbers: []int{11, 12},
	})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.DrinkServed{
		ID:          markDrinksServedID,
		MenuNumbers: []int{11, 12},
	})
}

func (suite *TabAggregateTestSuite) TestCannotServeUnorderedDrinks() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
		{MenuItem: 12, Description: "water", Price: 1.0},
	}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		ID:          markDrinksServedID,
		MenuNumbers: []int{11, 13},
	})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot serve drinks that were not ordered: [13]", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func (suite *TabAggregateTestSuite) TestCannotServeTheSameOrderedDrinkTwice() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		ID:          markDrinksServedID,
		MenuNumbers: []int{11},
	})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot serve drinks that were not ordered: [11]", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func (suite *TabAggregateTestSuite) TestCanCloseTabWhenPayingExactAmount() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
		{MenuItem: 12, Description: "water", Price: 1.0},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11, 12}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.CloseTab{ID: closeTabID, AmountPaid: 2.5})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.TabClosed{
		ID:          closeTabID,
		AmountPaid:  2.5,
		OrderAmount: 2.5,
		Tip:         0,
	})
}

func (suite *TabAggregateTestSuite) TestCanCloseTabWithTipWhenPayingMoreThanOrdered() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.CloseTab{ID: closeTabID, AmountPaid: 2.5})

	// Then
	assert.NoError(t, err)
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 1)
	suite.eventEmitter.AssertCalled(t, "EmitEvent", events.TabClosed{
		ID:          closeTabID,
		AmountPaid:  2.5,
		OrderAmount: 1.5,
		Tip:         1,
	})
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabWhenPayingLessThanOrdered() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.CloseTab{ID: closeTabID, AmountPaid: 1.0})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "not enough to cover tab, total served cost is: 1.5, but paid: 1", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabWithUnservedItems() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
		{MenuItem: 12, Description: "water", Price: 1.0},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11}})

	// When
	err := suite.tabAggregate.HandleCommand(commands.CloseTab{ID: closeTabID, AmountPaid: 1.5})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot close a tab with unserved items", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabTwice() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{ID: drinksOrderedEventID, Items: []domain.OrderedItem{
		{MenuItem: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinkServed{ID: drinksServedID, MenuNumbers: []int{11}})
	_ = suite.tabAggregate.ApplyEvent(events.TabClosed{ID: drinksServedID, AmountPaid: 1.5, OrderAmount: 1.5, Tip: 0})

	err := suite.tabAggregate.HandleCommand(commands.CloseTab{ID: closeTabID, AmountPaid: 1.5})

	assert.Error(t, err)
	assert.Equal(t, "cannot close a tab that is not open", err.Error())
	suite.eventEmitter.AssertNumberOfCalls(t, "EmitEvent", 0)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TabAggregateTestSuite))
}
