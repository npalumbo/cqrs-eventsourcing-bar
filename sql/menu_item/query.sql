-- name: GetAllMenuItems :many
SELECT id, description, price FROM menu_item ORDER BY id;

-- name: GetMenuItemsByIDs :many
SELECT id, description, price FROM menu_item WHERE id = ANY($1::int[]);
