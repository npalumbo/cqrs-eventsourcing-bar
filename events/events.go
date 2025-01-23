package events

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/shared"
	"reflect"

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

func GetEventTypeAsString(event Event) string {
	return reflect.TypeOf(event).Name()
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

type DrinksServed struct {
	BaseEvent
	MenuNumbers []int `json:"menu_numbers"`
}

type TabClosed struct {
	BaseEvent
	AmountPaid  float64 `json:"amount_paid"`
	OrderAmount float64 `json:"order_amount"`
	Tip         float64 `json:"tip"`
}

func UnmarshallPayload(typeName string, payload []byte) (Event, error) {
	switch typeName {
	case "TabOpened":
		var event TabOpened
		if err := json.Unmarshal(payload, &event); err != nil {
			return TabOpened{}, fmt.Errorf("could not create TabOpened event from payload: %s", payload)
		}
		return event, nil
	case "DrinksOrdered":
		var event DrinksOrdered
		if err := json.Unmarshal(payload, &event); err != nil {
			return DrinksOrdered{}, fmt.Errorf("could not create DrinksOrdered event from payload: %s", payload)
		}
		return event, nil
	case "DrinksServed":
		var event DrinksServed
		if err := json.Unmarshal(payload, &event); err != nil {
			return DrinksServed{}, fmt.Errorf("could not create DrinkServed event from payload: %s", payload)
		}
		return event, nil
	case "TabClosed":
		var event TabClosed
		if err := json.Unmarshal(payload, &event); err != nil {
			return TabClosed{}, fmt.Errorf("could not create TabClosed event from payload: %s", payload)
		}
		return event, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}
}
