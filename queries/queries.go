package queries

import (
	"errors"
	"fmt"
	"golangsevillabar/events"
	"slices"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
)

var todoByTab map[ksuid.KSUID]Tab = make(map[ksuid.KSUID]Tab)
var lock sync.RWMutex

type OpenTabQueries interface {
	ActiveTableNumbers() []int
	InvoiceForTable(table int) TabInvoice
	TabIdForTable(table int) ksuid.KSUID
	TabForTable(table int) TabStatus
	TodoListForWaiter(waiter string) map[int][]TabItem
}
type QueryEventHandler[E events.Event] interface {
	handle(e E) error
}

type tabOpenedEventHandler struct{}
type drinksOrderedEventHandler struct{}
type drinksServedEventHandler struct{}
type tabClosedEventHandler struct{}

func (t tabOpenedEventHandler) handle(e events.TabOpened) error {
	defer lock.Unlock()
	lock.Lock()
	todoByTab[e.ID] = Tab{
		TableNumber: e.TableNumber,
		Waiter:      e.Waiter,
		ToServe:     []TabItem{},
		Served:      []TabItem{},
	}
	return nil
}
func (d drinksOrderedEventHandler) handle(e events.DrinksOrdered) error {
	defer lock.Unlock()
	lock.Lock()
	tab := todoByTab[e.ID]
	addToServe := []TabItem{}
	for _, orderedItem := range e.Items {
		tabItem := TabItem{
			MenuNumber:  orderedItem.MenuItem,
			Description: orderedItem.Description,
			Price:       orderedItem.Price,
		}
		addToServe = append(addToServe, tabItem)
	}
	slices.AppendSeq(tab.ToServe, slices.Values(addToServe))
	return nil
}
func (d drinksServedEventHandler) handle(e events.DrinkServed) error {
	defer lock.Unlock()
	lock.Lock()
	tab := todoByTab[e.ID]
	for _, menuNumber := range e.MenuNumbers {
		foundElemWithMenuNumber := funk.Find(tab.ToServe, func(tabItem TabItem) bool {
			return tabItem.MenuNumber == menuNumber
		})
		tabItemFound, ok := foundElemWithMenuNumber.(TabItem)
		if ok {
			index := funk.IndexOf(tab.ToServe, tabItemFound)
			if index > -1 {
				tab.ToServe = slices.Delete(tab.ToServe, index, index)
				tab.Served = append(tab.Served, tabItemFound)
			}

		} else {
			return errors.New("found element that could not be transformed to TabItem")
		}
	}
	return nil
}
func (t tabClosedEventHandler) handle(e events.TabClosed) error {
	defer lock.Unlock()
	lock.Lock()
	delete(todoByTab, e.ID)
	return nil
}

type OpenTabs struct {
	tabOpenedEventHandler     QueryEventHandler[events.TabOpened]
	drinksOrderedEventHandler QueryEventHandler[events.DrinksOrdered]
	drinksServedEventHandler  QueryEventHandler[events.DrinkServed]
	tabClosedEventHandler     QueryEventHandler[events.TabClosed]
}

func (o OpenTabs) HandleEvent(e events.Event) error {
	switch event := e.(type) {
	case events.TabOpened:
		return o.tabOpenedEventHandler.handle(event)
	case events.DrinksOrdered:
		return o.drinksOrderedEventHandler.handle(event)
	case events.DrinkServed:
		return o.drinksServedEventHandler.handle(event)
	case events.TabClosed:
		return o.tabClosedEventHandler.handle(event)
	default:
		return fmt.Errorf("unexpected events.Event: %#v", e)
	}
}

func CreateOpenTabs() OpenTabs {
	return OpenTabs{
		tabOpenedEventHandler:     tabOpenedEventHandler{},
		drinksOrderedEventHandler: drinksOrderedEventHandler{},
		drinksServedEventHandler:  drinksServedEventHandler{},
		tabClosedEventHandler:     tabClosedEventHandler{},
	}
}

type TabInvoice struct {
	TabID            ksuid.KSUID
	TableNumber      int
	Items            []TabItem
	Total            float64
	HasUnservedItems bool
}

type TabStatus struct {
	TabID       ksuid.KSUID
	TableNumber int
	ToServe     []TabItem
	Served      []TabItem
}

type TabItem struct {
	MenuNumber  int
	Description string
	Price       float64
}

type Tab struct {
	TableNumber int
	Waiter      string
	ToServe     []TabItem
	Served      []TabItem
}
