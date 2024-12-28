package main

import (
	"golangsevillabar/commands"
	"net/http"
	"strconv"
	"strings"

	"github.com/segmentio/ksuid"
)

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

	id := q.Get("id")
	if id == "" {
		http.Error(w, "id needs to be defined", http.StatusBadRequest)
		return
	}

	items, err := parseToInts(q.Get("items"))

	if err != nil {
		http.Error(w, "could not parse items", http.StatusBadRequest)
		return
	}

	orderedItems, err := menuItemRepository.ReadItems(items)

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
