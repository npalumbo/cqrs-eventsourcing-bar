CREATE TABLE events (
    aggregate_id VARCHAR(28),
    sequence_number INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_type VARCHAR(512) NOT NULL,
    payload JSONB,
    PRIMARY KEY (aggregate_id, sequence_number)
);

CREATE TABLE menu_item (
    id INT NOT NULL,
    description VARCHAR(512) NOT NULL,
    price double precision,
    PRIMARY KEY (id)
);
