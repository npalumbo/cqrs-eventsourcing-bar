package ui

import (
	"golangsevillabar/queries"
	"slices"
)

type tabItemWithAmount struct {
	tabItem  queries.TabItem
	amount   int
	subTotal float64
}

func getTabItemsWithAmount(tabItems []queries.TabItem) []tabItemWithAmount {
	var amountsPerItemID map[int]int = make(map[int]int)
	var tabItemsByMenuNumbers map[int]queries.TabItem = make(map[int]queries.TabItem)

	for _, tabItem := range tabItems {
		_, ok := tabItemsByMenuNumbers[tabItem.MenuNumber]
		if !ok {
			tabItemsByMenuNumbers[tabItem.MenuNumber] = tabItem
		}
		amount, ok := amountsPerItemID[tabItem.MenuNumber]
		if !ok {
			amountsPerItemID[tabItem.MenuNumber] = 1
		} else {
			amountsPerItemID[tabItem.MenuNumber] = amount + 1
		}
	}

	tabItemWithAmounts := []tabItemWithAmount{}
	tabItemMenuNumbers := make([]int, 0, len(amountsPerItemID))
	for k := range amountsPerItemID {
		tabItemMenuNumbers = append(tabItemMenuNumbers, k)
	}
	slices.Sort(tabItemMenuNumbers)

	for _, menuNumber := range tabItemMenuNumbers {
		amount := amountsPerItemID[menuNumber]
		tabItemWithAmounts = append(tabItemWithAmounts, tabItemWithAmount{
			amount:   amount,
			tabItem:  tabItemsByMenuNumbers[menuNumber],
			subTotal: float64(amount) * tabItemsByMenuNumbers[menuNumber].Price,
		})
	}

	return tabItemWithAmounts
}
