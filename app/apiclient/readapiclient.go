package apiclient

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/readservice/model"
	"io"
	"net/http"
)

type ReadClient struct {
	httpClient *http.Client
	url        string
}

func NewReadClient(httpClient *http.Client, url string) *ReadClient {
	return &ReadClient{httpClient: httpClient, url: url}
}

func (c *ReadClient) GetActiveTables() (model.ActiveTableNumbersResponse, error) {
	response := model.ActiveTableNumbersResponse{}
	req, err := http.NewRequest("GET", c.url+"/activeTableNumbers", nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *ReadClient) GetTabIdForTable(tableNumber int) (model.TabIdForTableResponse, error) {
	response := model.TabIdForTableResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tabIdForTable?table_number=%d", c.url, tableNumber), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *ReadClient) GetTabForTable(tableNumber int) (model.TabForTableResponse, error) {
	response := model.TabForTableResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tabForTable?table_number=%d", c.url, tableNumber), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *ReadClient) GetInvoiceForTable(tableNumber int) (model.InvoiceForTableResponse, error) {
	response := model.InvoiceForTableResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/invoiceForTable?table_number=%d", c.url, tableNumber), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *ReadClient) GetTodoListForWaiter(waiter string) (model.TodoListForWaiterResponse, error) {
	response := model.TodoListForWaiterResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/todoListForWaiter?waiter=%s", c.url, waiter), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *ReadClient) GetAllMenuItems() (model.AllMenuItemsResponse, error) {
	response := model.AllMenuItemsResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/allMenuItems", c.url), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func processResponse[T any](c *ReadClient, req *http.Request, response T) (T, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return response, err
	}

	return response, nil
}
