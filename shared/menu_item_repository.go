package shared

import "context"

//go:generate mockery --name MenuItemRepository
type MenuItemRepository interface {
	ReadItems(ctx context.Context, menuItems []int) ([]MenuItem, error)
	ReadAllItems(ctx context.Context) ([]MenuItem, error)
}

type MenuItem struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
