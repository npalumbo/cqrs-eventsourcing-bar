package messaging_test

import (
	"golangsevillabar/events"
	"golangsevillabar/messaging"
	"testing"
	"time"

	mock_events "golangsevillabar/events/mocks"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NatsRoundtripTestSuite struct {
	suite.Suite
	NatsEventEmitter    *messaging.NatsEventEmitter
	NatsEventSubscriber *messaging.NatsEventSubscriber
	natsServer          *server.Server
	eventListener       *mock_events.EventListener
}

func (suite *NatsRoundtripTestSuite) SetupSuite() {
	eventListener := mock_events.NewEventListener(suite.T())
	opts := server.Options{
		Host: "localhost",
		Port: 1234, // Let the server choose a random port
	}
	suite.natsServer = server.New(&opts)

	suite.natsServer.Start()

	natsURL := suite.natsServer.ClientURL()
	natsEventEmitter, err := messaging.NewNatsEventEmitter(natsURL)
	if err != nil {
		suite.T().Error(err.Error())
	}
	natsEventSubscriber, err := messaging.NewNatsEventSubscriber(natsURL, eventListener)
	if err != nil {
		suite.T().Error(err.Error())
	}
	suite.NatsEventEmitter = natsEventEmitter
	suite.NatsEventSubscriber = natsEventSubscriber
	suite.eventListener = eventListener
}

func (suite *NatsRoundtripTestSuite) TestEmitEventShouldCallEventListener() {
	// Given
	err := suite.NatsEventSubscriber.OnCreatedEvent()

	if err != nil {
		assert.Fail(suite.T(), err.Error())
	}

	aggregateId := ksuid.New()
	tabOpened := events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: aggregateId},
		TableNumber: 1,
		Waiter:      "waiter_1",
	}
	suite.eventListener.On("HandleEvent", tabOpened).Return(nil)

	// When
	err = suite.NatsEventEmitter.EmitEvent(tabOpened)
	if err != nil {
		assert.Fail(suite.T(), err.Error())
	}

	// Then
	time.Sleep(20 * time.Millisecond)
	suite.eventListener.AssertExpectations(suite.T())
	suite.eventListener.AssertNumberOfCalls(suite.T(), "HandleEvent", 1)
}

func (suite *NatsRoundtripTestSuite) TearDownSuite() {
	suite.natsServer.Shutdown()
}

func TestNatsTestSuite(t *testing.T) {
	suite.Run(t, new(NatsRoundtripTestSuite))
}
