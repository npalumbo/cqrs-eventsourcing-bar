package events

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/shared"

	"github.com/segmentio/ksuid"
)

type Event interface {
	GetID() ksuid.KSUID
}

type BaseEvent struct {
	ID ksuid.KSUID `json:"id"`
}

func (event BaseEvent) GetID() ksuid.KSUID {
	return event.ID
}

type TabOpened struct {
	BaseEvent
	TableNumber int    `json:"table_number"`
	Waiter      string `json:"waiter"`
}

type DrinksOrdered struct {
	BaseEvent
	Items []shared.OrderedItem `json:"items"`
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
		var event DrinksOrdered
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return DrinksOrdered{}, fmt.Errorf("could not create DrinksOrdered event from payload: %s", payload)
		}
		return event, nil
	case "DrinkServed":
		var event DrinkServed
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return DrinkServed{}, fmt.Errorf("could not create DrinkServed event from payload: %s", payload)
		}
		return event, nil
	case "TabClosed":
		var event TabClosed
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			return TabClosed{}, fmt.Errorf("could not create TabClosed event from payload: %s", payload)
		}
		return event, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}
}
