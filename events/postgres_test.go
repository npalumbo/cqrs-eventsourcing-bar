package events_test

import (
	"context"
	"log"
	"reflect"
	"testing"

	"golangsevillabar/events"
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
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
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

	aggregateId, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")
	events, err := suite.eventStorePostgres.LoadEvents(aggregateId)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Equal(t, aggregateId, events[0].GetID())
	assert.Equal(t, "TabOpened", reflect.TypeOf(events[0]).String())

}

func TestPostgresEventStoreTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresEventStoreTestSuite))
}
