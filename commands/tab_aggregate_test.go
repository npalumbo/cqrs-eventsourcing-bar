package commands_test

import (
	"golangsevillabar/commands"
	"golangsevillabar/events"
	"golangsevillabar/shared"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TabAggregateTestSuite struct {
	suite.Suite
	tabAggregate commands.Aggregate
}

func (suite *TabAggregateTestSuite) SetupTest() {
	suite.tabAggregate = commands.TabAggregateFactory{}.CreateAggregate()
}

func (suite *TabAggregateTestSuite) TestCanOpenTab() {

	commandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.OpenTab{
		BaseCommand: commands.BaseCommand{ID: commandID},
		TableNumber: 0,
		Waiter:      "waiter_1",
	})

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: commandID},
		TableNumber: 0,
		Waiter:      "waiter_1",
	}}, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanNotOrderWithUnOpenedTab() {
	commandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.PlaceOrder{
		BaseCommand: commands.BaseCommand{ID: commandID},
		Items:       []shared.MenuItem{{ID: 11, Description: "beer", Price: 1.5}},
	})

	// Then
	assert.Error(t, err)
	assert.Empty(t, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanOrderWhenTabIsOpen() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	placeOrderCommandID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	// Given(events.TabOpened{ID: tabOpenedEventID, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.PlaceOrder{
		BaseCommand: commands.BaseCommand{ID: placeOrderCommandID},
		Items:       []shared.MenuItem{{ID: 11, Description: "beer", Price: 1.5}},
	})

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: placeOrderCommandID},
		Items:     []shared.MenuItem{{ID: 11, Description: "beer", Price: 1.5}},
	}}, newEvents)
}

func (suite *TabAggregateTestSuite) TestOrderedDrinksCanBeServed() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
		{ID: 12, Description: "water", Price: 1.0},
	}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		BaseCommand: commands.BaseCommand{ID: markDrinksServedID},
		MenuNumbers: []int{11, 12},
	})

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{events.DrinksServed{
		BaseEvent:   events.BaseEvent{ID: markDrinksServedID},
		MenuNumbers: []int{11, 12},
	}}, newEvents)
}

func (suite *TabAggregateTestSuite) TestCannotServeUnorderedDrinks() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
		{ID: 12, Description: "water", Price: 1.0},
	}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		BaseCommand: commands.BaseCommand{ID: markDrinksServedID},
		MenuNumbers: []int{11, 13},
	})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot serve drinks that were not ordered: [13]", err.Error())
	assert.Empty(t, newEvents)
}

func (suite *TabAggregateTestSuite) TestCannotServeTheSameOrderedDrinkTwice() {
	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	markDrinksServedID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.MarkDrinksServed{
		BaseCommand: commands.BaseCommand{ID: markDrinksServedID},
		MenuNumbers: []int{11},
	})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot serve drinks that were not ordered: [11]", err.Error())
	assert.Empty(t, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanCloseTabWhenPayingExactAmount() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
		{ID: 12, Description: "water", Price: 1.0},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11, 12}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.CloseTab{BaseCommand: commands.BaseCommand{ID: closeTabID}, AmountPaid: 2.5})

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{events.TabClosed{
		BaseEvent:   events.BaseEvent{ID: closeTabID},
		AmountPaid:  2.5,
		OrderAmount: 2.5,
		Tip:         0,
	}}, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanCloseTabWithTipWhenPayingMoreThanOrdered() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.CloseTab{BaseCommand: commands.BaseCommand{ID: closeTabID}, AmountPaid: 2.5})

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{events.TabClosed{
		BaseEvent:   events.BaseEvent{ID: closeTabID},
		AmountPaid:  2.5,
		OrderAmount: 1.5,
		Tip:         1,
	}}, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabWhenPayingLessThanOrdered() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.CloseTab{BaseCommand: commands.BaseCommand{ID: closeTabID}, AmountPaid: 1.0})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "not enough to cover tab, total served cost is: 1.5, but paid: 1", err.Error())
	assert.Empty(t, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabWithUnservedItems() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
		{ID: 12, Description: "water", Price: 1.0},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11}})

	// When
	newEvents, err := suite.tabAggregate.HandleCommand(commands.CloseTab{BaseCommand: commands.BaseCommand{ID: closeTabID}, AmountPaid: 1.5})

	// Then
	assert.Error(t, err)
	assert.Equal(t, "cannot close a tab with unserved items", err.Error())
	assert.Empty(t, newEvents)
}

func (suite *TabAggregateTestSuite) TestCanotCloseTabTwice() {

	tabOpenedEventID, _ := ksuid.NewRandom()
	drinksOrderedEventID, _ := ksuid.NewRandom()
	drinksServedID, _ := ksuid.NewRandom()
	closeTabID, _ := ksuid.NewRandom()
	t := suite.T()

	// Given
	_ = suite.tabAggregate.ApplyEvent(events.TabOpened{BaseEvent: events.BaseEvent{ID: tabOpenedEventID}, Waiter: "waiter_1", TableNumber: 0})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksOrdered{BaseEvent: events.BaseEvent{ID: drinksOrderedEventID}, Items: []shared.MenuItem{
		{ID: 11, Description: "beer", Price: 1.5},
	}})
	_ = suite.tabAggregate.ApplyEvent(events.DrinksServed{BaseEvent: events.BaseEvent{ID: drinksServedID}, MenuNumbers: []int{11}})
	_ = suite.tabAggregate.ApplyEvent(events.TabClosed{BaseEvent: events.BaseEvent{ID: drinksServedID}, AmountPaid: 1.5, OrderAmount: 1.5, Tip: 0})

	newEvents, err := suite.tabAggregate.HandleCommand(commands.CloseTab{BaseCommand: commands.BaseCommand{ID: closeTabID}, AmountPaid: 1.5})

	assert.Error(t, err)
	assert.Equal(t, "cannot close a tab that is not open", err.Error())
	assert.Empty(t, newEvents)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TabAggregateTestSuite))
}
