package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/segmentio/ksuid"
)

type postgresEventStore struct {
	conn *pgx.Conn
}

func (es *postgresEventStore) LoadEvents(aggregateID ksuid.KSUID) ([]Event, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	rows, err := es.conn.Query(ctx, "SELECT event_type, payload FROM events WHERE aggregate_id = $1 ORDER BY sequence_number ASC", aggregateID.String())
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

func (es *postgresEventStore) SaveEvents(aggregateID ksuid.KSUID, lastKnownEventID int, events []Event) error {
	ctx := context.Background()
	tx, err := es.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert new events
	for i, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO events (aggregate_id, sequence_number, timestamp, event_type, payload) VALUES ($1, $2, NOW(), $3)", aggregateID, lastKnownEventID+i, payload)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func NewPostgresEventStore(ctx context.Context, connStr string) (EventStore, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		slog.Error("Unable to connect to database", slog.String("error", err.Error()))
		return nil, err
	}
	return &postgresEventStore{
		conn: conn,
	}, nil
}
