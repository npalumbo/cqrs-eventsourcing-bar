package main

import (
	"context"
	"cqrseventsourcingbar/events"
	"cqrseventsourcingbar/messaging"
	"cqrseventsourcingbar/queries"
	"cqrseventsourcingbar/readservice/service"
	"cqrseventsourcingbar/shared"
	"fmt"
)

func main() {
	openTabQueries := queries.CreateOpenTabs()
	natsEventSubscriber, err := messaging.NewNatsEventSubscriber("nats://localhost:4222", openTabQueries)
	panicIfErrors(err)

	ctx := context.Background()

	const dbConnectionString = "postgresql://postgres:mysecretpassword@localhost:5432/mydb"
	eventStore, err := events.NewPostgresEventStore(ctx, dbConnectionString)
	panicIfErrors(err)

	events, err := eventStore.LoadAllEvents(ctx)
	panicIfErrors(err)

	for _, event := range events {
		err = openTabQueries.HandleEvent(event)
		panicIfErrors(err)
	}

	menuItemRepository, err := shared.NewPostgresMenuItemRepository(ctx, dbConnectionString)
	panicIfErrors(err)

	err = natsEventSubscriber.OnCreatedEvent()
	panicIfErrors(err)

	readService := service.CreateReadService(8081, openTabQueries, menuItemRepository)

	err = readService.Start()
	panicIfErrors(err)
}

func panicIfErrors(err error) {
	if err != nil {
		panic(fmt.Sprintf("error: %s, not starting app", err.Error()))
	}
}
