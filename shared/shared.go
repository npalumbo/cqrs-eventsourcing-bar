package shared

type OrderedItem struct {
	MenuItem    int     `json:"menu_item"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
