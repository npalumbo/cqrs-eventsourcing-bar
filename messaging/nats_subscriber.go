package messaging

import (
	"bytes"
	"cqrseventsourcingbar/events"
	"encoding/gob"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type NatsEventSubscriber struct {
	conn            *nats.Conn
	eventCreatedSub *nats.Subscription
	eventListener   events.EventListener
}

func (n *NatsEventSubscriber) OnCreatedEvent() error {
	msg := wrappedEvent{}
	var err error
	n.eventCreatedSub, err = n.conn.Subscribe("event", func(m *nats.Msg) {
		err := n.decodeMessage(m.Data, &msg)
		if err != nil {
			logIncomingEventError(err)
			return
		}

		event, err := events.UnmarshallPayload(msg.EventType, msg.Payload)
		if err != nil {
			logIncomingEventError(err)
			return
		}

		err = n.eventListener.HandleEvent(event)

		if err != nil {
			logIncomingEventError(err)
			return
		}
	})
	return err
}

func NewNatsEventSubscriber(url string, eventListener events.EventListener) (*NatsEventSubscriber, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventSubscriber{
		conn:          conn,
		eventListener: eventListener,
	}, nil
}

func (n *NatsEventSubscriber) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
	if n.eventCreatedSub != nil {
		err := n.eventCreatedSub.Unsubscribe()
		if err != nil {
			slog.Error("error closing subscription", slog.Any("error", err.Error()))
		}
	}
}

func (n *NatsEventSubscriber) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func logIncomingEventError(err error) {
	slog.Error("error when processing incoming event", slog.Any("error", err.Error()))
}
