package events

import (
	"context"

	"github.com/segmentio/ksuid"
)

//go:generate mockery --name EventStore
type EventStore interface {
	LoadEvents(ctx context.Context, aggregateID ksuid.KSUID) ([]Event, error)
	LoadAllEvents(ctx context.Context) ([]Event, error)
	SaveEvents(ctx context.Context, aggregateID ksuid.KSUID, previousEventCount int, events []Event) error
}
