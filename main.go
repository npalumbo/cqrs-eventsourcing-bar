package main

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/events"

	"github.com/segmentio/ksuid"
)

func main() {
	id, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")

	to := events.TabOpened{
		BaseEvent:   events.BaseEvent{ID: id},
		TableNumber: 1,
		Waiter:      "w1",
	}

	data, _ := json.Marshal(to)

	fmt.Println(string(data))
}
