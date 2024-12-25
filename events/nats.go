package events

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type NatsEmitterSubscriber struct {
	conn            *nats.Conn
	eventCreatedSub *nats.Subscription
	eventListener   EventListener
}

type wrappedEvent struct {
	EventType string
	Payload   []byte
}

func NewNatsEmitterSubscriber(url string, eventListener EventListener) (*NatsEmitterSubscriber, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEmitterSubscriber{
		conn:          conn,
		eventListener: eventListener,
	}, nil
}

func (n *NatsEmitterSubscriber) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
	if n.eventCreatedSub != nil {
		n.eventCreatedSub.Unsubscribe()
	}
}

func (n *NatsEmitterSubscriber) encodeMessage(event Event) ([]byte, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	wrappedEvent := wrappedEvent{
		EventType: GetEventTypeAsString(event),
		Payload:   payload,
	}
	b := bytes.Buffer{}
	err = gob.NewEncoder(&b).Encode(wrappedEvent)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *NatsEmitterSubscriber) EmitEvent(event Event) error {
	data, err := n.encodeMessage(event)
	if err != nil {
		return err
	}
	return n.conn.Publish("event", data)
}

func (n *NatsEmitterSubscriber) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEmitterSubscriber) OnCreatedEvent() (err error) {
	msg := wrappedEvent{}
	n.eventCreatedSub, err = n.conn.Subscribe("event", func(m *nats.Msg) {
		err = n.decodeMessage(m.Data, &msg)
		if err == nil {
			event, err := UnmarshallPayload(msg.EventType, msg.Payload)
			if err == nil {
				n.eventListener.HandleEvent(event)
			}
		}
		if err != nil {
			slog.Error("error when processing incoming event", slog.Any("error", err.Error()))
		}
	})
	return err

}
