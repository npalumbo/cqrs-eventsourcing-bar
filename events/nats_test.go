package events_test

import (
	"golangsevillabar/events"
	"testing"
	"time"

	mock_events "golangsevillabar/events/mocks"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NatsTestSuite struct {
	suite.Suite
	natsEmitterSusbcriber *events.NatsEmitterSubscriber
	natsServer            *server.Server
	eventListener         *mock_events.EventListener
}

func (suite *NatsTestSuite) SetupSuite() {
	eventListener := mock_events.NewEventListener(suite.T())
	opts := server.Options{
		Host: "localhost",
		Port: 1234, // Let the server choose a random port
	}
	suite.natsServer = server.New(&opts)

	suite.natsServer.Start()

	natsURL := suite.natsServer.ClientURL()
	natsEmitterSusbcriber, err := events.NewNatsEmitterSubscriber(natsURL, eventListener)
	if err != nil {
		suite.T().Fail()
	}
	suite.natsEmitterSusbcriber = natsEmitterSusbcriber
	suite.eventListener = eventListener
}

func (suite *NatsTestSuite) TestEmitEvent() {
	// Given
	suite.natsEmitterSusbcriber.OnCreatedEvent()
	aggregateId := ksuid.New()
	tabOpened := events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: aggregateId},
		TableNumber: 1,
		Waiter:      "waiter_1",
	}
	suite.eventListener.On("HandleEvent", tabOpened).Return(nil)
	// When

	err := suite.natsEmitterSusbcriber.EmitEvent(tabOpened)
	if err != nil {
		assert.Fail(suite.T(), "error emitting event")
	}

	// Then
	time.Sleep(20 * time.Millisecond)
	suite.eventListener.AssertExpectations(suite.T())
	suite.eventListener.AssertNumberOfCalls(suite.T(), "HandleEvent", 1)
}

func (suite *NatsTestSuite) TearDownSuite() {
	suite.natsServer.Shutdown()
}

func TestNatsTestSuite(t *testing.T) {
	suite.Run(t, new(NatsTestSuite))
}
