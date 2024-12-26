package messaging

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"golangsevillabar/events"

	"github.com/nats-io/nats.go"
)

type NatsEventEmitter struct {
	conn *nats.Conn
}

type wrappedEvent struct {
	EventType string
	Payload   []byte
}

func NewNatsEventEmitter(url string) (*NatsEventEmitter, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventEmitter{
		conn: conn,
	}, nil
}

func (n *NatsEventEmitter) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

func (n *NatsEventEmitter) encodeMessage(event events.Event) ([]byte, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	wrappedEvent := wrappedEvent{
		EventType: events.GetEventTypeAsString(event),
		Payload:   payload,
	}
	b := bytes.Buffer{}
	err = gob.NewEncoder(&b).Encode(wrappedEvent)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *NatsEventEmitter) EmitEvent(event events.Event) error {
	data, err := n.encodeMessage(event)
	if err != nil {
		return err
	}
	return n.conn.Publish("event", data)
}
