package queries_test

import (
	"golangsevillabar/domain"
	"golangsevillabar/events"
	"golangsevillabar/queries"
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
	suite.openTabQueries.HandleEvent(events.TabOpened{
		ID:          tabId,
		TableNumber: 1,
		Waiter:      "Charles",
	})

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId, invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 0.0, invoice.Total)
	assert.Empty(suite.T(), invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabForTable.TabID, tabId)
	assert.Equal(suite.T(), tabForTable.TableNumber, 1)
	assert.Equal(suite.T(), tabForTable.ToServe, []queries.TabItem{})

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {}}, todoListForWaiter)
}

func (suite *QueriesTestSuite) TestAnOpenTabWithOneOrder() {
	tabId := ksuid.New()
	suite.openTabQueries.HandleEvent(events.TabOpened{
		ID:          tabId,
		TableNumber: 1,
		Waiter:      "Charles",
	})

	suite.openTabQueries.HandleEvent(events.DrinksOrdered{
		ID: tabId,
		Items: []domain.OrderedItem{{
			MenuItem:    10,
			Description: "Water",
			Price:       1,
		}},
	})

	activeTableNumbers := suite.openTabQueries.ActiveTableNumbers()
	assert.Equal(suite.T(), []int{1}, activeTableNumbers)

	invoice, err := suite.openTabQueries.InvoiceForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, invoice.HasUnservedItems)
	assert.Equal(suite.T(), tabId, invoice.TabID)
	assert.Equal(suite.T(), 1, invoice.TableNumber)
	assert.Equal(suite.T(), 0.0, invoice.Total)
	assert.Empty(suite.T(), invoice.Items)

	tabForTable, err := suite.openTabQueries.TabForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabForTable.TabID)
	assert.Equal(suite.T(), 1, tabForTable.TableNumber)
	assert.Equal(suite.T(), []queries.TabItem{{MenuNumber: 10, Description: "Water", Price: 1}}, tabForTable.ToServe)

	tabIdForTable1, err := suite.openTabQueries.TabIdForTable(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tabId, tabIdForTable1)

	todoListForWaiter := suite.openTabQueries.TodoListForWaiter("Charles")
	assert.Equal(suite.T(), map[int][]queries.TabItem{1: {{MenuNumber: 10, Description: "Water", Price: 1}}}, todoListForWaiter)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(QueriesTestSuite))
}
