package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/events"
	"golangsevillabar/messaging"
	"net/http"

	"github.com/segmentio/ksuid"
)

var dispatcher *commands.Dispatcher

type OpenTabRequest struct {
	TableNumber int    `json:"table_number"`
	Waiter      string `json:"waiter"`
}

func openTabHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request OpenTabRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = dispatcher.DispatchCommand(r.Context(), commands.OpenTab{
		BaseCommand: commands.BaseCommand{ID: ksuid.New()},
		TableNumber: request.TableNumber,
		Waiter:      request.Waiter,
	})

	if err != nil {
		http.Error(w, "Error processing openTab request", http.StatusInternalServerError)
		return
	}
}

func main() {
	ctx := context.Background()
	eventStore, err := events.NewPostgresEventStore(ctx, "")

	panicIfErrors(err)

	eventEmitter, err := messaging.NewNatsEventEmitter("")

	panicIfErrors(err)

	dispatcher = commands.CreateCommandDispatcher(eventStore, eventEmitter, commands.TabAggregateFactory{})

	http.HandleFunc("/openTab", openTabHandler)

	fmt.Println("Write server listening on :8080")
	err = http.ListenAndServe(":8080", nil)

	panicIfErrors(err)
}

func panicIfErrors(err error) {
	if err != nil {
		panic(fmt.Sprintf("error: %s, not starting app", err.Error()))
	}
}
