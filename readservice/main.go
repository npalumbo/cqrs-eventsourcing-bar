package main

import (
	"fmt"
	"golangsevillabar/messaging"
	"golangsevillabar/queries"
	"golangsevillabar/readservice/service"
)

func main() {
	openTabQueries := queries.CreateOpenTabs()
	natsEventSubscriber, err := messaging.NewNatsEventSubscriber("nats://localhost:4222", openTabQueries)
	panicIfErrors(err)

	err = natsEventSubscriber.OnCreatedEvent()
	panicIfErrors(err)

	readService := service.CreateReadService(8081, openTabQueries)

	err = readService.Start()
	panicIfErrors(err)
}

func panicIfErrors(err error) {
	if err != nil {
		panic(fmt.Sprintf("error: %s, not starting app", err.Error()))
	}
}
