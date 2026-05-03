package events

import (
	"context"
	"encoding/json"
	"log/slog"

	eventsdb "cqrseventsourcingbar/internal/db/events"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/ksuid"
)

type postgresEventStore struct {
	queries *eventsdb.Queries
	conn    *pgx.Conn
}

func (es *postgresEventStore) LoadAllEvents(ctx context.Context) ([]Event, error) {
	rows, err := es.queries.LoadAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, row := range rows {
		event, err := UnmarshallPayload(row.EventType, row.Payload)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (es *postgresEventStore) LoadEvents(ctx context.Context, aggregateID ksuid.KSUID) ([]Event, error) {
	rows, err := es.queries.LoadEventsByAggregateID(ctx, aggregateID.String())
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, row := range rows {
		event, err := UnmarshallPayload(row.EventType, row.Payload)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (es *postgresEventStore) SaveEvents(ctx context.Context, aggregateID ksuid.KSUID, previousEventCount int, eventsToSave []Event) (err error) {
	tx, err := es.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				slog.Error("error, rollback", slog.String("error", rbErr.Error()))
			}
		}
	}()

	txQueries := eventsdb.New(tx)

	for i, event := range eventsToSave {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = txQueries.SaveEvent(ctx, eventsdb.SaveEventParams{
			AggregateID:    aggregateID.String(),
			SequenceNumber: int32(previousEventCount + i + 1),
			EventType:      GetEventTypeAsString(event),
			Payload:        payload,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func NewPostgresEventStore(ctx context.Context, connStr string) (EventStore, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		slog.Error("unable to connect to database", slog.String("error", err.Error()))
		return nil, err
	}
	return &postgresEventStore{
		queries: eventsdb.New(conn),
		conn:    conn,
	}, nil
}
