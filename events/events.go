package events

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/shared"

	"github.com/segmentio/ksuid"
)

type Event interface{}

type BaseEvent struct {
	ID ksuid.KSUID
}

func (event BaseEvent) GetID() ksuid.KSUID {
	return event.ID
}

type TabOpened struct {
	BaseEvent
	TableNumber int
	Waiter      string
}

type DrinksOrdered struct {
	BaseEvent
	Items []shared.OrderedItem
}

type DrinkServed struct {
	BaseEvent
	MenuNumbers []int
}

type TabClosed struct {
	BaseEvent
	AmountPaid  float64
	OrderAmount float64
	Tip         float64
}

func UnmarshallPayload(typeName, payload string) (Event, error) {
	switch typeName {
	case "TabOpened":
		var event TabOpened
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return TabOpened{}, fmt.Errorf("could not create TabOpened event from payload: %s", payload)
		}
		return event, nil
	case "DrinksOrdered":
		var event TabOpened
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return TabOpened{}, fmt.Errorf("could not create TabOpened event from payload: %s", payload)
		}
		return event, nil
	case "DrinkServed":
		var event TabOpened
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return TabOpened{}, fmt.Errorf("could not create TabOpened event from payload: %s", payload)
		}
		return event, nil
	case "TabClosed":
		var event TabOpened
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return TabOpened{}, fmt.Errorf("could not create TabOpened event from payload: %s", payload)
		}
		return event, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}
}
