package shared

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	menuitemdb "cqrseventsourcingbar/internal/db/menu_item"

	"github.com/jackc/pgx/v5"
)

type postgresMenuItemRepository struct {
	queries *menuitemdb.Queries
}

func (p *postgresMenuItemRepository) ReadAllItems(ctx context.Context) ([]MenuItem, error) {
	dbItems, err := p.queries.GetAllMenuItems(ctx)
	if err != nil {
		return nil, err
	}

	allItems := make([]MenuItem, len(dbItems))
	for i, item := range dbItems {
		allItems[i] = MenuItem{
			ID:          int(item.ID),
			Description: item.Description,
			Price:       item.Price.Float64,
		}
	}

	return allItems, nil
}

func (p *postgresMenuItemRepository) ReadItems(ctx context.Context, menuItems []int) ([]MenuItem, error) {
	slices.Sort(menuItems)
	originalItems := slices.Clone(menuItems)
	uniqueItems := slices.Compact(menuItems)

	dbIDs := make([]int32, len(uniqueItems))
	for i, id := range uniqueItems {
		dbIDs[i] = int32(id)
	}

	dbItems, err := p.queries.GetMenuItemsByIDs(ctx, dbIDs)
	if err != nil {
		return nil, err
	}

	retrievedItems := make(map[int]MenuItem)
	for _, item := range dbItems {
		retrievedItems[int(item.ID)] = MenuItem{
			ID:          int(item.ID),
			Description: item.Description,
			Price:       item.Price.Float64,
		}
	}

	amountOfRetrievedItems := len(retrievedItems)
	amountOfInputItems := len(uniqueItems)
	if amountOfRetrievedItems != amountOfInputItems {
		return nil, fmt.Errorf("requested %d distinct items, but read from DB %d distinct items", amountOfInputItems, amountOfRetrievedItems)
	}

	orderedItems := make([]MenuItem, len(originalItems))
	for i, id := range originalItems {
		orderedItems[i] = retrievedItems[id]
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
		queries: menuitemdb.New(conn),
	}, nil
}
