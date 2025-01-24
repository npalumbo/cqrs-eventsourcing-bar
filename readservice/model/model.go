package model

import "golangsevillabar/queries"

type QueryResponse[T any] struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Data  T      `json:"data"`
}

type ActiveTableNumbersResponse QueryResponse[[]int]

type TabIdForTableResponse QueryResponse[string]

type TabForTableResponse QueryResponse[queries.TabStatus]

type InvoiceForTableResponse QueryResponse[queries.TabInvoice]

type TodoListForWaiter QueryResponse[map[int][]queries.TabItem]

// type ActiveTableNumbersResponse struct {
// 	ActiveTables []int `json:"active_tables"`
// }

// type TabIdForTableResponse struct {
// 	TabId string `json:"tab_id"`
// }

// type TabForTableResponse struct {
// 	TabStatus queries.TabStatus `json:"tab_status"`
// }

// type InvoiceForTableResponse struct {
// 	TabInvoice queries.TabInvoice `json:"tab_invoice"`
// }

// type TodoListForWaiterResponse struct {
// 	TabItemsForTable map[int][]queries.TabItem `json:"tab_items_for_table"`
// }
