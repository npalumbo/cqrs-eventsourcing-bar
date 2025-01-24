package events_test

import (
	"context"
	"log"
	"testing"

	"golangsevillabar/events"
	"golangsevillabar/shared"
	"golangsevillabar/testhelpers"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresEventStoreTestSuite struct {
	suite.Suite
	pgContainer        *testhelpers.PostgresContainer
	eventStorePostgres events.EventStore
	ctx                context.Context
}

func (suite *PostgresEventStoreTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.T(), suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	eventStorePostgres, err := events.NewPostgresEventStore(suite.ctx, suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.eventStorePostgres = eventStorePostgres
}

func (suite *PostgresEventStoreTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *PostgresEventStoreTestSuite) TestLoadEvents() {
	t := suite.T()
	// Given
	aggregateId, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")
	// When
	loadedEvents, err := suite.eventStorePostgres.LoadEvents(context.TODO(), aggregateId)
	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, loadedEvents)
	assert.Len(t, loadedEvents, 2)
	assert.Equal(t, aggregateId, loadedEvents[0].GetID())
	tabOpened, ok := loadedEvents[0].(events.TabOpened)
	assert.True(t, ok)
	assert.Equal(t, events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: aggregateId},
		TableNumber: 1,
		Waiter:      "w1",
	}, tabOpened)
	drinksOrdered, ok := loadedEvents[1].(events.DrinksOrdered)
	assert.True(t, ok)
	assert.Equal(t, events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: aggregateId},
		Items: []shared.MenuItem{{
			ID:          1,
			Description: "water",
			Price:       1.5,
		}},
	}, drinksOrdered)

}

func (suite *PostgresEventStoreTestSuite) TestLoadAllEvents() {
	t := suite.T()
	// Given
	aggregateId1, _ := ksuid.Parse("1qPTBJCN6ib7iJ6WaIVvoSmySSV")
	aggregateId2, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")
	// When
	loadedEvents, err := suite.eventStorePostgres.LoadAllEvents(context.TODO())
	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, loadedEvents)
	assert.Len(t, loadedEvents, 3)
	assert.Equal(t, aggregateId1, loadedEvents[0].GetID())
	tabOpened, ok := loadedEvents[0].(events.TabOpened)
	assert.True(t, ok)
	assert.Equal(t, events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: aggregateId1},
		TableNumber: 2,
		Waiter:      "w2",
	}, tabOpened)
	tabOpened2, ok := loadedEvents[1].(events.TabOpened)
	assert.True(t, ok)
	assert.Equal(t, events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: aggregateId2},
		TableNumber: 1,
		Waiter:      "w1",
	}, tabOpened2)

}

func (suite *PostgresEventStoreTestSuite) TestSaveEventsErrorsIfWeAttemptToOverrideExistingEvent() {
	// Given
	aggregateId, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")

	eventsToSave := []events.Event{events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: aggregateId},
		Items: []shared.MenuItem{{
			ID:          1,
			Description: "water",
			Price:       1.5,
		}}}}

	// When
	err := suite.eventStorePostgres.SaveEvents(context.TODO(), aggregateId, 1, eventsToSave)
	// Then
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "ERROR: duplicate key value violates unique constraint \"events_pkey\" (SQLSTATE 23505)", err.Error())
}

func (suite *PostgresEventStoreTestSuite) TestSaveEvents() {
	// Given
	aggregateId, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")

	eventsToSave := []events.Event{events.DrinksOrdered{
		BaseEvent: events.BaseEvent{ID: aggregateId},
		Items: []shared.MenuItem{{
			ID:          1,
			Description: "water",
			Price:       1.5,
		}}}}

	// When
	err := suite.eventStorePostgres.SaveEvents(context.TODO(), aggregateId, 2, eventsToSave)
	// Then
	assert.NoError(suite.T(), err)
}

func TestPostgresEventStoreTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresEventStoreTestSuite))
}
