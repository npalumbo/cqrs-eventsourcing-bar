-- name: LoadAllEvents :many
SELECT event_type, payload FROM events ORDER BY aggregate_id, sequence_number ASC;

-- name: LoadEventsByAggregateID :many
SELECT event_type, payload FROM events WHERE aggregate_id = $1 ORDER BY sequence_number ASC;

-- name: SaveEvent :exec
INSERT INTO events (aggregate_id, sequence_number, timestamp, event_type, payload)
VALUES ($1, $2, NOW(), $3, $4);
