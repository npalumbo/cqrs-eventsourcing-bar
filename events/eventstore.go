package events

import "github.com/segmentio/ksuid"

//go:generate mockery --name EventStore
type EventStore interface {
	LoadEvents(aggregateID ksuid.KSUID) ([]Event, error)
	SaveEvents(aggregateID ksuid.KSUID, lastKnownEventID int, events []Event) error
}
