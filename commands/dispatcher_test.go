package commands_test

import (
	"context"
	"errors"
	"fmt"
	"golangsevillabar/commands"
	mock_commands "golangsevillabar/commands/mocks"
	"golangsevillabar/events"
	mock_events "golangsevillabar/events/mocks"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DispatcherTestSuite struct {
	suite.Suite
	dispatcher   *commands.Dispatcher
	eventEmitter *mock_events.EventEmitter
	eventStore   *mock_events.EventStore
	aggregate    *mock_commands.Aggregate
	ctx          context.Context
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
	suite.ctx = context.TODO()
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToLoadEvents() {
	// Given
	aggregateId := ksuid.New()
	errorLoadingEvents := errors.New("all broken")
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return(nil, errorLoadingEvents)

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
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
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(errors.New("all broken"))

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error applying past event [BaseEvent-#0] for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToHandleCommand() {
	// Given
	aggregateId := ksuid.New()
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(nil)
	suite.aggregate.On("HandleCommand", commands.BaseCommand{ID: aggregateId}).Return(nil, errors.New("all broken"))

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error handling command [BaseCommand] for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToSaveEvents() {
	// Given
	aggregateId := ksuid.New()
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(nil)
	suite.aggregate.On("HandleCommand", commands.BaseCommand{ID: aggregateId}).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.eventStore.On("SaveEvents", suite.ctx, aggregateId, 1, []events.Event{events.BaseEvent{ID: aggregateId}}).Return(errors.New("all broken"))

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error when saving events for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func (suite *DispatcherTestSuite) TestDispatcherReturnsErrorWhenFailingToEmitEvents() {
	// Given
	aggregateId := ksuid.New()
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(nil)
	suite.aggregate.On("HandleCommand", commands.BaseCommand{ID: aggregateId}).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.eventStore.On("SaveEvents", suite.ctx, aggregateId, 1, []events.Event{events.BaseEvent{ID: aggregateId}}).Return(nil)
	suite.eventEmitter.On("EmitEvent", events.BaseEvent{ID: aggregateId}).Return(errors.New("all broken"))

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
		commands.BaseCommand{aggregateId},
	)

	// Then
	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), fmt.Sprintf("error when emitting event [BaseEvent] for aggregate: %s, reason: all broken", aggregateId), err.Error())
	}
}

func (suite *DispatcherTestSuite) TestDispatcher() {
	// Given
	aggregateId := ksuid.New()
	suite.eventStore.On("LoadEvents", suite.ctx, aggregateId).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.aggregate.On("ApplyEvent", events.BaseEvent{ID: aggregateId}).Return(nil)
	suite.aggregate.On("HandleCommand", commands.BaseCommand{ID: aggregateId}).Return([]events.Event{events.BaseEvent{ID: aggregateId}}, nil)
	suite.eventStore.On("SaveEvents", suite.ctx, aggregateId, 1, []events.Event{events.BaseEvent{ID: aggregateId}}).Return(nil)
	suite.eventEmitter.On("EmitEvent", events.BaseEvent{ID: aggregateId}).Return(nil)

	// When
	err := suite.dispatcher.DispatchCommand(
		context.TODO(),
		commands.BaseCommand{aggregateId},
	)

	// Then
	assert.NoError(suite.T(), err)
}

func TestDispatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}
