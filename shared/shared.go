package shared

type MenuItemRepository interface {
	ReadItems(menuItems []int) ([]OrderedItem, error)
}

type OrderedItem struct {
	MenuItem    int     `json:"menu_item"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
