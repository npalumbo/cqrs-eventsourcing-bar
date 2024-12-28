package main

import (
	"context"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/events"
	"golangsevillabar/messaging"
	"golangsevillabar/shared"
	"net/http"
)

var dispatcher *commands.Dispatcher
var menuItemRepository shared.MenuItemRepository

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
