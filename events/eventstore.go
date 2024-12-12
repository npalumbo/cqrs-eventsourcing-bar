package events

import "github.com/segmentio/ksuid"

type EventStore interface {
	LoadEvents(aggregateID ksuid.KSUID) ([]Event, error)
	SaveEvents(aggregateID ksuid.KSUID, lastKnownEventID int, events []Event) error
}
