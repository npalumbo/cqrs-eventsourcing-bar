CREATE TABLE events (
    aggregate_id VARCHAR(20) PRIMARY KEY,
    sequence_number SERIAL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_type VARCHAR(128),
    payload JSONB
);