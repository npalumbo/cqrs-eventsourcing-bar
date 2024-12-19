CREATE TABLE events (
    aggregate_id VARCHAR(28),
    sequence_number INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_type VARCHAR(512) NOT NULL,
    payload JSONB,
    PRIMARY KEY (aggregate_id, sequence_number)
);

INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":1,"waiter":"w1"}');
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 2, 'DrinksOrdered', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","items":[{"menu_item":1,"description":"water","price":1.5}]}');
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('1qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"1qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":2,"waiter":"w2"}');