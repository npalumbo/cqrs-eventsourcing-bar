package main

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/writeservice/model"
	"io"
	"log/slog"
	"net/http"

	"github.com/segmentio/ksuid"
)

func setupServer() error {
	http.HandleFunc("/openTab", openTabHandler)
	http.HandleFunc("/placeOrder", placeOrderHandler)
	http.HandleFunc("/markDrinksServed", markDrinksServedHandler)
	http.HandleFunc("/closeTab", closeTabHandler)

	slog.Info("Write server listening on :8080")

	return http.ListenAndServe(":8080", nil)
}

func openTabHandler(w http.ResponseWriter, r *http.Request) {
	var request model.OpenTabRequest
	shouldReturn := readRequest(w, r, &request)
	if shouldReturn {
		return
	}

	err := dispatcher.DispatchCommand(r.Context(), commands.OpenTab{
		BaseCommand: commands.BaseCommand{ID: ksuid.New()},
		TableNumber: request.TableNumber,
		Waiter:      request.Waiter,
	})

	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing openTab request: %v", err), http.StatusInternalServerError)
		return
	}

	returnJsonOk(w)
}

func placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	var request model.PlaceOrderRequest
	shouldReturn := readRequest(w, r, &request)
	if shouldReturn {
		return
	}

	id, err := ksuid.Parse(request.TabId)

	if err != nil {
		returnJsonError(w, "could not parse id", http.StatusBadRequest)
		return
	}

	orderedItems, err := menuItemRepository.ReadItems(r.Context(), request.MenuItems)

	if err != nil {
		returnJsonError(w, "could not read items from DB", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.PlaceOrder{
		BaseCommand: commands.BaseCommand{ID: id},
		Items:       orderedItems,
	})

	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing placeOrder request: %v", err), http.StatusInternalServerError)
		return
	}

	returnJsonOk(w)
}

func markDrinksServedHandler(w http.ResponseWriter, r *http.Request) {
	var request model.MarkDrinksServedRequest
	shouldReturn := readRequest(w, r, &request)
	if shouldReturn {
		return
	}

	id, err := ksuid.Parse(request.TabId)

	if err != nil {
		returnJsonError(w, "could not parse id", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.MarkDrinksServed{
		BaseCommand: commands.BaseCommand{ID: id},
		MenuNumbers: request.MenuNumbers,
	})

	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing markDrinksServed request: %v", err), http.StatusInternalServerError)
		return
	}

	returnJsonOk(w)
}

func closeTabHandler(w http.ResponseWriter, r *http.Request) {
	var request model.CloseTabRequest
	shouldReturn := readRequest(w, r, &request)
	if shouldReturn {
		return
	}

	id, err := ksuid.Parse(request.TabId)

	if err != nil {
		returnJsonError(w, "could not parse id", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.CloseTab{
		BaseCommand: commands.BaseCommand{ID: id},
		AmountPaid:  request.AmountPaid,
	})

	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing closeTab request: %v", err), http.StatusInternalServerError)
		return
	}

	returnJsonOk(w)
}

func readRequest[T any](w http.ResponseWriter, r *http.Request, data *T) (errored bool) {
	if r.Method != http.MethodPost {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		errored = true
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		returnJsonError(w, "Error reading request body", http.StatusBadRequest)
		errored = true
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		returnJsonError(w, "Invalid JSON request", http.StatusBadRequest)
		errored = true
	}
	return
}

func returnJsonError(w http.ResponseWriter, error string, code int) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	response := model.CommandReponse{
		OK:    false,
		Error: error,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding json, original error: %s", error), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing json response, original error: %s", error), http.StatusInternalServerError)
	}
}

func returnJsonOk(w http.ResponseWriter) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	response := model.CommandReponse{
		OK:    true,
		Error: "",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error encoding json, command processed sucesssfully", http.StatusOK)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "error writing json response, command processed sucesssfully", http.StatusOK)
	}
}
