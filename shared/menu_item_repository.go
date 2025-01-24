package shared

import "context"

//go:generate mockery --name MenuItemRepository
type MenuItemRepository interface {
	ReadItems(ctx context.Context, menuItems []int) ([]OrderedItem, error)
	ReadAllItems(ctx context.Context) ([]OrderedItem, error)
}

type OrderedItem struct {
	MenuItem    int     `json:"menu_item"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
