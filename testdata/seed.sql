INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":1,"waiter":"w1"}') ON CONFLICT DO NOTHING;
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('2qPTBJCN6ib7iJ6WaIVvoSmySSV', 2, 'DrinksOrdered', '{"id":"2qPTBJCN6ib7iJ6WaIVvoSmySSV","items":[{"id":1,"description":"water","price":1.5}]}') ON CONFLICT DO NOTHING;
INSERT INTO events(aggregate_id, sequence_number, event_type, payload) VALUES ('1qPTBJCN6ib7iJ6WaIVvoSmySSV', 1, 'TabOpened', '{"id":"1qPTBJCN6ib7iJ6WaIVvoSmySSV","table_number":2,"waiter":"w2"}') ON CONFLICT DO NOTHING;

INSERT INTO menu_item(id, description, price) VALUES (1, 'blue water', 1.0) ON CONFLICT DO NOTHING;
INSERT INTO menu_item(id, description, price) VALUES (2, 'red water', 2.0) ON CONFLICT DO NOTHING;
INSERT INTO menu_item(id, description, price) VALUES (3, 'green water', 3.0) ON CONFLICT DO NOTHING;
