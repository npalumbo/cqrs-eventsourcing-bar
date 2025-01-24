package shared

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/jackc/pgx/v5"
)

type postgresMenuItemRepository struct {
	conn *pgx.Conn
}

func (p *postgresMenuItemRepository) ReadAllItems(ctx context.Context) ([]MenuItem, error) {
	rows, err := p.conn.Query(ctx, "SELECT id, description, price FROM menu_item ORDER by id")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	allItems := []MenuItem{}

	for rows.Next() {
		var id int
		var description string
		var price float64
		if err := rows.Scan(&id, &description, &price); err != nil {
			return nil, err
		}
		allItems = append(allItems, MenuItem{
			ID:          id,
			Description: description,
			Price:       price,
		})
	}

	return allItems, nil
}

func (p *postgresMenuItemRepository) ReadItems(ctx context.Context, menuItems []int) ([]MenuItem, error) {
	slices.Sort(menuItems)
	originalItems := slices.Clone(menuItems)
	uniqueItems := slices.Compact(menuItems)
	rows, err := p.conn.Query(ctx, "SELECT id, description, price FROM menu_item WHERE id = any($1)", uniqueItems)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var retrievedItems map[int]MenuItem = make(map[int]MenuItem)
	for rows.Next() {
		var id int
		var description string
		var price float64
		if err := rows.Scan(&id, &description, &price); err != nil {
			return nil, err
		}
		retrievedItems[id] = MenuItem{
			ID:          id,
			Description: description,
			Price:       price,
		}
	}

	amountOfRetrievedItems := len(retrievedItems)
	amountOfInputItems := len(uniqueItems)
	if amountOfRetrievedItems != amountOfInputItems {
		return nil, fmt.Errorf("requested %d distinct items, but read from DB %d distinct items", amountOfInputItems, amountOfRetrievedItems)
	}

	orderedItems := []MenuItem{}

	for _, i := range originalItems {
		orderedItems = append(orderedItems, retrievedItems[i])
	}

	return orderedItems, nil
}

func NewPostgresMenuItemRepository(ctx context.Context, connStr string) (MenuItemRepository, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		slog.Error("unable to connect to database", slog.String("error", err.Error()))
		return nil, err
	}
	return &postgresMenuItemRepository{
		conn: conn,
	}, nil
}
