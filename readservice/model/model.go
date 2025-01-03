package model

import "golangsevillabar/queries"

type ActiveTableNumbersResponse struct {
	ActiveTables []int  `json:"active_tables"`
	OK           bool   `json:"ok"`
	Error        string `json:"error"`
}

type InvoiceForTableRequest struct {
	TableNumber int `json:"table_number"`
}

type InvoiceForTableResponse struct {
	TabInvoice queries.TabInvoice `json:"tab_invoice"`
	OK         bool               `json:"ok"`
	Error      string             `json:"error"`
}

type TabIdForTableRequest struct {
	TableNumber int `json:"table_number"`
}

type TabIdForTableResponse struct {
	TabId string `json:"tab_id"`
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type TabForTableRequest struct {
	TableNumber int `json:"table_number"`
}

type TabForTableResponse struct {
	TabId string `json:"tab_id"`
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type TodoListForWaiterRequest struct {
	Waiter string `json:"waiter"`
}

type TodoListForWaiterResponse struct {
	TabItems []queries.TabItem `json:"tab_items"`
	OK       bool              `json:"ok"`
	Error    string            `json:"error"`
}
