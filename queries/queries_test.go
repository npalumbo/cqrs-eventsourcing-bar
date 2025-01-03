package queries_test

import (
	"golangsevillabar/events"
	"golangsevillabar/queries"
	"golangsevillabar/shared"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueriesTestSuite struct {
	suite.Suite
	openTabQueries queries.OpenTabQueries
}

func (suite *QueriesTestSuite) SetupTest() {
	suite.openTabQueries = queries.CreateOpenTabs()
}

func (suite *QueriesTestSuite) TestNoOpenTabs() {
	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Empty(suite.T(), activeTableNumbers)

	_, err := suite.openTabQueries.InvoiceForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	_, err = suite.openTabQueries.TabForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	_, err = suite.openTabQueries.TabIdForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Jenkins")
	assert.Empty(suite.T(), todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAnOpenTab() {
	tabId := ksuid.New()
	err := suite.openTabQueries.HandleEvent(events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: tabId},
		TableNumber: 1,
		Waiter:      "Charles",
	})
	assert.NoError(suite.T(), err)

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId.String(), invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 0.0, invoice.Total)
	assert.Empty(suite.T(), invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId.String(), tabForTable.TabID)
	assert.Equal(suite.T(), 1, tabForTable.TableNumber)
	assert.Equal(suite.T(), []queries.TabItem{}, tabForTable.ToServe)

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {}}, todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAnOpenTabWithOneOrder() {
	tabId := ksuid.New()
	err := suite.openTabQueries.HandleEvent(events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: tabId},
		TableNumber: 1,
		Waiter:      "Charles",
	})

	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: tabId},
		Items: []shared.OrderedItem{{
			MenuItem:    10,
			Description: "Water",
			Price:       1,
		}},
	})

	assert.NoError(suite.T(), err)

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId.String(), invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 0.0, invoice.Total)
	assert.Empty(suite.T(), invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId.String(), tabForTable.TabID)
	assert.Equal(suite.T(), 1, tabForTable.TableNumber)
	assert.Equal(suite.T(), []queries.TabItem{{MenuNumber: 10, Description: "Water", Price: 1}}, tabForTable.ToServe)

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {{MenuNumber: 10, Description: "Water", Price: 1}}}, todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAnOpenTabWithTwoOrdersOnlyOneServed() {
	tabId := ksuid.New()
	err := suite.openTabQueries.HandleEvent(events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: tabId},
		TableNumber: 1,
		Waiter:      "Charles",
	})

	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: tabId},
		Items: []shared.OrderedItem{{
			MenuItem:    10,
			Description: "Water",
			Price:       1,
		}},
	})

	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: tabId},
		Items: []shared.OrderedItem{{
			MenuItem:    11,
			Description: "Beer",
			Price:       2,
		}},
	})

	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinkServed{
		BaseEvent:   events.BaseEvent{ID: tabId},
		MenuNumbers: []int{10},
	})

	assert.NoError(suite.T(), err)

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId.String(), invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 1.0, invoice.Total)
	assert.Equal(suite.T(), []queries.TabItem{{MenuNumber: 10, Description: "Water", Price: 1}}, invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId.String(), tabForTable.TabID)
	assert.Equal(suite.T(), 1, tabForTable.TableNumber)
	assert.Equal(suite.T(), []queries.TabItem{{MenuNumber: 11, Description: "Beer", Price: 2}}, tabForTable.ToServe)

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {{MenuNumber: 11, Description: "Beer", Price: 2}}}, todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAnOpenTabWithTwoOrdersBothServed() {
	tabId := ksuid.New()
	err := suite.openTabQueries.HandleEvent(events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: tabId},
		TableNumber: 1,
		Waiter:      "Charles",
	})
	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: tabId},
		Items: []shared.OrderedItem{{
			MenuItem:    10,
			Description: "Water",
			Price:       1,
		}},
	})
	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: tabId},
		Items: []shared.OrderedItem{{
			MenuItem:    11,
			Description: "Beer",
			Price:       2,
		}},
	})
	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.DrinkServed{
		BaseEvent:   events.BaseEvent{ID: tabId},
		MenuNumbers: []int{10},
	})
	assert.NoError(suite.T(), err)
	err = suite.openTabQueries.HandleEvent(events.DrinkServed{
		BaseEvent:   events.BaseEvent{ID: tabId},
		MenuNumbers: []int{11},
	})
	assert.NoError(suite.T(), err)

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId.String(), invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 3.0, invoice.Total)
	assert.Equal(suite.T(), []queries.TabItem{{MenuNumber: 10, Description: "Water", Price: 1}, {
		MenuNumber:  11,
		Description: "Beer",
		Price:       2,
	}}, invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId.String(), tabForTable.TabID)
	assert.Equal(suite.T(), 1, tabForTable.TableNumber)
	assert.Empty(suite.T(), tabForTable.ToServe)

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {}}, todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAfterCloseThereIsNoData() {
	tabId := ksuid.New()
	err := suite.openTabQueries.HandleEvent(events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: tabId},
		TableNumber: 1,
		Waiter:      "Charles",
	})
	assert.NoError(suite.T(), err)

	err = suite.openTabQueries.HandleEvent(events.TabClosed{
		BaseEvent:   events.BaseEvent{ID: tabId},
		AmountPaid:  0.0,
		OrderAmount: 0.0,
		Tip:         0.0,
	})
	assert.NoError(suite.T(), err)

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Empty(suite.T(), activeTableNumbers)

	_, err = suite.openTabQueries.InvoiceForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	_, err = suite.openTabQueries.TabForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	_, err = suite.openTabQueries.TabIdForTable(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "couldn't find a tab for table: 1", err.Error())

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Jenkins")
	assert.Empty(suite.T(), todoListForWaiter)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(QueriesTestSuite))
}
