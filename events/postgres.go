package events

import (
	"database/sql"

	"github.com/segmentio/ksuid"
)

type postgresEventStore struct {
	db *sql.DB
}

func (es *postgresEventStore) LoadEvents(aggregateID ksuid.KSUID) ([]Event, error) {
	rows, err := es.db.Query("SELECT event_type, payload FROM events WHERE aggregate_id = $1 ORDER BY sequence_number ASC", aggregateID.String())
	// rows, err := es.db.QueryContext(ctx, "SELECT event_type, payload FROM events WHERE aggregate_id = $1 ORDER BY sequence_number ASC", aggregateID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var eventType string
		var payload []byte
		if err := rows.Scan(&eventType, &payload); err != nil {
			return nil, err
		}
		event, err := UnmarshallPayload(eventType, string(payload))
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (p *postgresEventStore) SaveEvents(aggregateID ksuid.KSUID, lastKnownEventID int, events []Event) error {
	panic("unimplemented")
}

func NewPostgresEventStore(db *sql.DB) EventStore {
	return &postgresEventStore{db: db}
}
