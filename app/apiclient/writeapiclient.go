package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golangsevillabar/writeservice/model"
	"io"
	"log/slog"
	"net/http"
)

type WriteClient struct {
	httpClient *http.Client
	url        string
}

func NewWriteClient(httpClient *http.Client, url string) *WriteClient {
	return &WriteClient{httpClient: httpClient, url: url}
}

func (w *WriteClient) ExecuteCommand(commandRequest interface{}) error {
	uri, err := resolveUri(commandRequest)
	if err != nil {
		slog.Error("Error resolving uri from request:", slog.Any("error", err))
		return err
	}
	jsonData, err := json.Marshal(commandRequest)
	if err != nil {
		slog.Error("Error marshalling JSON:", slog.Any("error", err))
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", w.url, uri), bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating http request:", slog.Any("error", err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	commandResponse := model.CommandReponse{}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		slog.Error("Error on transport:", slog.Any("error", err))
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", slog.Any("error", err))
		return err
	}

	err = json.Unmarshal(body, &commandResponse)
	if err != nil {
		slog.Error("Error parsing response JSON:", slog.Any("error", err))
		return err
	}

	if commandResponse.OK {
		return nil
	}

	return errors.New(commandResponse.Error)
}

func resolveUri(request interface{}) (string, error) {
	switch request.(type) {
	case model.OpenTabRequest:
		return "openTab", nil
	case model.PlaceOrderRequest:
		return "placeOrder", nil
	case model.MarkDrinksServedRequest:
		return "markDrinksServed", nil
	case model.CloseTabRequest:
		return "closeTab", nil
	default:
		return "", errors.New("unsupported request type")
	}
}
