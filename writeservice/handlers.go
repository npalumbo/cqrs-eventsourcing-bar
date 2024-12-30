package main

import (
	"fmt"
	"golangsevillabar/commands"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	tableNumberStr := q.Get("table_number")
	tableNumber, err := strconv.Atoi(tableNumberStr)
	if err != nil {
		http.Error(w, "Invalid table number", http.StatusBadRequest)
		return
	}

	waiter := q.Get("waiter")
	if waiter == "" {
		http.Error(w, "Waiter needs to be defined", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.OpenTab{
		BaseCommand: commands.BaseCommand{ID: ksuid.New()},
		TableNumber: tableNumber,
		Waiter:      waiter,
	})

	if err != nil {
		http.Error(w, "Error processing openTab request", http.StatusInternalServerError)
		return
	}
}

func placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	idStr := q.Get("id")
	if idStr == "" {
		http.Error(w, "id needs to be defined", http.StatusBadRequest)
		return
	}

	id, err := ksuid.Parse(idStr)

	if err != nil {
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	items, err := parseToInts(q.Get("items"))

	if err != nil {
		http.Error(w, "could not parse items", http.StatusBadRequest)
		return
	}

	orderedItems, err := menuItemRepository.ReadItems(r.Context(), items)

	if err != nil {
		http.Error(w, "could read items from DB", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.PlaceOrder{
		BaseCommand: commands.BaseCommand{ID: id},
		Items:       orderedItems,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing placeOrder request: %v", err), http.StatusInternalServerError)
		return
	}

}

func markDrinksServedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	idStr := q.Get("id")
	if idStr == "" {
		http.Error(w, "id needs to be defined", http.StatusBadRequest)
		return
	}

	id, err := ksuid.Parse(idStr)

	if err != nil {
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	menuNumbers, err := parseToInts(q.Get("menu_numbers"))

	if err != nil {
		http.Error(w, "could not parse items", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.MarkDrinksServed{
		BaseCommand: commands.BaseCommand{ID: id},
		MenuNumbers: menuNumbers,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing markDrinksServed request: %v", err), http.StatusInternalServerError)
		return
	}
}

func closeTabHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	idStr := q.Get("id")
	if idStr == "" {
		http.Error(w, "id needs to be defined", http.StatusBadRequest)
		return
	}

	id, err := ksuid.Parse(idStr)

	if err != nil {
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	amountPaid, err := strconv.ParseFloat(q.Get("amount_paid"), 64)

	if err != nil {
		http.Error(w, "could not parse amount paid", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.CloseTab{
		BaseCommand: commands.BaseCommand{ID: id},
		AmountPaid:  amountPaid,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing closeTab request: %v", err), http.StatusInternalServerError)
		return
	}
}

func parseToInts(str string) ([]int, error) {
	var ints []int
	for _, s := range strings.Split(str, ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}
