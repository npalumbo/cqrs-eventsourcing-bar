CREATE TABLE events (
    aggregate_id VARCHAR(28),
    sequence_number INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_type VARCHAR(512) NOT NULL,
    payload JSONB,
    PRIMARY KEY (aggregate_id, sequence_number)
);

INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":1,"waiter":"waiter 1"}');
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 2, 'DrinksOrdered', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","items":[{"id":1,"description":"blue water","price":1.0},{"id":2,"description":"red water","price":2.0}]}');
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 3, 'DrinksServed', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","menu_numbers":[1,2]}');
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('1qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"1qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":2,"waiter":"waiter 2"}');

CREATE TABLE menu_item (
    id INT NOT NULL,
    description VARCHAR(512) NOT NULL,
    price double precision,
    PRIMARY KEY (id)
);

INSERT INTO menu_item(id, description, price) VALUES (1, 'blue water', 1.0);
INSERT INTO menu_item(id, description, price) VALUES (2, 'red water', 2.0);
INSERT INTO menu_item(id, description, price) VALUES (3, 'green water', 3.0);