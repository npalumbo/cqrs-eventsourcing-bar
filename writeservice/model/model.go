package model

type OpenTabRequest struct {
	TableNumber int    `json:"table_number"`
	Waiter      string `json:"waiter"`
}

type PlaceOrderRequest struct {
	TabId     string `json:"tab_id"`
	MenuItems []int  `json:"menu_items"`
}

type MarkDrinksServedRequest struct {
	TabId       string `json:"tab_id"`
	MenuNumbers []int  `json:"menu_numbers"`
}

type CloseTabRequest struct {
	TabId      string  `json:"tab_id"`
	AmountPaid float64 `json:"amount_paid"`
}

type CommandReponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}
