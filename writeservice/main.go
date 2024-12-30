package main

import (
	"context"
	"fmt"
	"golangsevillabar/commands"
	"golangsevillabar/events"
	"golangsevillabar/messaging"
	"golangsevillabar/shared"
)

var dispatcher *commands.Dispatcher
var menuItemRepository shared.MenuItemRepository

func main() {
	ctx := context.Background()
	const dbConnectionString = "postgresql://postgres:mysecretpassword@localhost:5432/mydb"
	eventStore, err := events.NewPostgresEventStore(ctx, dbConnectionString)
	panicIfErrors(err)

	menuItemRepository, err = shared.NewPostgresMenuItemRepository(ctx, dbConnectionString)
	panicIfErrors(err)

	eventEmitter, err := messaging.NewNatsEventEmitter("nats://localhost:4222")

	panicIfErrors(err)

	dispatcher = commands.CreateCommandDispatcher(eventStore, eventEmitter, commands.TabAggregateFactory{})

	err = setupServer()

	panicIfErrors(err)
}

func panicIfErrors(err error) {
	if err != nil {
		panic(fmt.Sprintf("error: %s, not starting app", err.Error()))
	}
}
