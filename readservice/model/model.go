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

type TodoListForWaiterResponse QueryResponse[map[int][]queries.TabItem]
