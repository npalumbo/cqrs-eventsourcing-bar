package apiclient

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/readservice/model"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{httpClient: httpClient, url: url}
}

func (c *Client) GetActiveTables() (model.QueryResponse[model.ActiveTableNumbersResponse], error) {
	response := model.QueryResponse[model.ActiveTableNumbersResponse]{}
	req, err := http.NewRequest("GET", c.url+"/activeTableNumbers", nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func (c *Client) GetTabIdForTable(tableNumber int) (model.QueryResponse[model.TabIdForTableResponse], error) {
	response := model.QueryResponse[model.TabIdForTableResponse]{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tabIdForTable?table_number=%d", c.url, tableNumber), nil)
	if err != nil {
		return response, err
	}

	return processResponse(c, req, response)
}

func processResponse[T any](c *Client, req *http.Request, response T) (T, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return response, err
	}
	defer resp.Body.Close()

	// Read the response body
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
