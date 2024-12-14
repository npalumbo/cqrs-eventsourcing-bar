package commands_test

import (
	"errors"
	"fmt"
	"golangsevillabar/commands"
	mock_commands "golangsevillabar/commands/mocks"
	"golangsevillabar/events"
	mock_events "golangsevillabar/events/mocks"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DispatcherTestSuite struct {
	suite.Suite
	dispatcher   *commands.Dispatcher
	eventEmitter *mock_events.EventEmitter
	eventStore   *mock_events.EventStore
	aggregate    *mock_commands.Aggregate
}

func (suite *DispatcherTestSuite) SetupTest() {
	eventEmitter := mock_events.NewEventEmitter(suite.T())
	eventStore := mock_events.NewEventStore(suite.T())
	aggregateFactory := mock_commands.NewAggregateFactory(suite.T())
	aggregate := mock_commands.NewAggregate(suite.T())
	aggregateFactory.On("CreateAggregate").Return(aggregate)
	suite.dispatcher = commands.CreateCommandDispatcher(eventStore, eventEmitter, aggregateFactory)
	suite.eventEmitter = eventEmitter
	suite.eventStore = eventStore
	suite.aggregate = aggregate
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToLoadEvents() {
	// Given
	errorLoadingEvents := errors.New("all broken")
	suite.eventStore.On("LoadEvents", mock.Anything).Return(nil, errorLoadingEvents)

	// When
	aggregateId := ksuid.New()
	err := suite.dispatcher.DispatchCommand(
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error loading events for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToApplyEventOnAggregate() {
	// Given
	aggregateId := ksuid.New()
	suite.eventStore.On("LoadEvents", mock.Anything).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(errors.New("all broken"))

	// When
	err := suite.dispatcher.DispatchCommand(
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error applying past event [BaseEvent-#0] for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func TestDispatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}
