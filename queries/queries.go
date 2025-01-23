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

//go:generate mockery --name OpenTabQueries
type OpenTabQueries interface {
	ActiveTableNumbers() []int
	InvoiceForTable(table int) (TabInvoice, error)
	TabIdForTable(table int) (ksuid.KSUID, error)
	TabForTable(table int) (TabStatus, error)
	TodoListForWaiter(waiter string) map[int][]TabItem
	events.EventListener
}

func (o *openTabs) handleTabOpened(e events.TabOpened) error {
	defer o.lock.Unlock()
	o.lock.Lock()
	o.todoByTab[e.ID] = &Tab{
		TableNumber: e.TableNumber,
		Waiter:      e.Waiter,
		ToServe:     []TabItem{},
		Served:      []TabItem{},
	}
	return nil
}
func (o *openTabs) handleDrinksOrdered(e events.DrinksOrdered) error {
	defer o.lock.Unlock()
	o.lock.Lock()
	tab := o.todoByTab[e.ID]
	addToServe := []TabItem{}
	for _, orderedItem := range e.Items {
		tabItem := TabItem{
			MenuNumber:  orderedItem.MenuItem,
			Description: orderedItem.Description,
			Price:       orderedItem.Price,
		}
		addToServe = append(addToServe, tabItem)
	}
	tab.ToServe = slices.AppendSeq(tab.ToServe, slices.Values(addToServe))
	return nil
}
func (o *openTabs) handleDrinksServed(e events.DrinksServed) error {
	defer o.lock.Unlock()
	o.lock.Lock()
	tab := o.todoByTab[e.ID]
	for _, menuNumber := range e.MenuNumbers {
		foundElemWithMenuNumber := funk.Find(tab.ToServe, func(tabItem TabItem) bool {
			return tabItem.MenuNumber == menuNumber
		})
		tabItemFound, ok := foundElemWithMenuNumber.(TabItem)
		if ok {
			index := funk.IndexOf(tab.ToServe, tabItemFound)
			if index > -1 {
				tab.ToServe = slices.Delete(tab.ToServe, index, index+1)
				tab.Served = append(tab.Served, tabItemFound)
			}

		} else {
			return errors.New("found element that could not be transformed to TabItem")
		}
	}
	return nil
}
func (o *openTabs) handleTabClosed(e events.TabClosed) error {
	defer o.lock.Unlock()
	o.lock.Lock()
	delete(o.todoByTab, e.ID)
	return nil
}

type openTabs struct {
	todoByTab map[ksuid.KSUID]*Tab
	lock      sync.RWMutex
}

func (o *openTabs) ActiveTableNumbers() []int {
	defer o.lock.RUnlock()
	o.lock.RLock()
	tableNumbers := []int{}
	for _, todo := range o.todoByTab {
		tableNumbers = append(tableNumbers, todo.TableNumber)
	}
	slices.Sort(tableNumbers)
	return tableNumbers
}

func (o *openTabs) InvoiceForTable(table int) (TabInvoice, error) {
	defer o.lock.RUnlock()
	o.lock.RLock()

	tabId, err := o.TabIdForTable(table)

	if err != nil {
		return TabInvoice{}, err
	}

	tab := o.todoByTab[tabId]

	mapped := funk.Map(tab.Served, func(t TabItem) float64 {
		return t.Price
	})
	total := funk.SumFloat64(mapped.([]float64))

	return TabInvoice{
		TabID:            tabId.String(),
		TableNumber:      table,
		Items:            slices.Clone(tab.Served),
		Total:            total,
		HasUnservedItems: len(tab.ToServe) > 0,
	}, nil

}

func (o *openTabs) TabForTable(table int) (TabStatus, error) {
	defer o.lock.RUnlock()
	o.lock.RLock()

	tabId, err := o.TabIdForTable(table)

	if err != nil {
		return TabStatus{}, err
	}

	tab := o.todoByTab[tabId]

	return TabStatus{
		TabID:       tabId.String(),
		TableNumber: table,
		ToServe:     slices.Clone(tab.ToServe),
		Served:      slices.Clone(tab.Served),
	}, nil

}

func (o *openTabs) TabIdForTable(table int) (ksuid.KSUID, error) {
	defer o.lock.RUnlock()
	o.lock.RLock()
	tabId, _ := funk.FindKey(o.todoByTab, func(tab *Tab) bool {
		return tab.TableNumber == table
	})

	if tabIdAsKSUID, ok := tabId.(ksuid.KSUID); ok {
		return tabIdAsKSUID, nil
	}

	return ksuid.KSUID{}, fmt.Errorf("couldn't find a tab for table: %d", table)

}

func (o *openTabs) TodoListForWaiter(waiter string) map[int][]TabItem {
	todoListForWaiter := make(map[int][]TabItem)

	for _, v := range o.todoByTab {
		if v.Waiter == waiter {

			tabItems := []TabItem{}
			tabItems = append(tabItems, v.ToServe...)
			todoListForWaiter[v.TableNumber] = tabItems
		}
	}

	return todoListForWaiter
}

func (o *openTabs) HandleEvent(e events.Event) error {
	switch event := e.(type) {
	case events.TabOpened:
		return o.handleTabOpened(event)
	case events.DrinksOrdered:
		return o.handleDrinksOrdered(event)
	case events.DrinksServed:
		return o.handleDrinksServed(event)
	case events.TabClosed:
		return o.handleTabClosed(event)
	default:
		return fmt.Errorf("unexpected events.Event: %#v", e)
	}
}

func CreateOpenTabs() OpenTabQueries {
	return &openTabs{
		todoByTab: make(map[ksuid.KSUID]*Tab),
		lock:      sync.RWMutex{},
	}
}

type TabInvoice struct {
	TabID            string    `json:"tab_id"`
	TableNumber      int       `json:"table_number"`
	Items            []TabItem `json:"items"`
	Total            float64   `json:"total"`
	HasUnservedItems bool      `json:"has_unserved_items"`
}

type TabStatus struct {
	TabID       string    `json:"tab_id"`
	TableNumber int       `json:"table_number"`
	ToServe     []TabItem `json:"to_serve"`
	Served      []TabItem `json:"served"`
}

type TabItem struct {
	MenuNumber  int     `json:"menu_number"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Tab struct {
	TableNumber int       `json:"table_number"`
	Waiter      string    `json:"waiter"`
	ToServe     []TabItem `json:"to_serve"`
	Served      []TabItem `json:"served"`
}
