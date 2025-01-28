package ui

import (
	"fmt"
	"golangsevillabar/queries"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTabItemList(tabItemsWithAmount *[]tabItemWithAmount) *widget.List {

	return widget.NewList(func() int { return len(*tabItemsWithAmount) }, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, widget.NewLabel(""), widget.NewLabel(""), nil)
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		tabItemWithAmount := (*tabItemsWithAmount)[lii]
		borderContainer := co.(*fyne.Container)
		objects := borderContainer.Objects
		leftLabel := objects[0].(*widget.Label)
		rightLabel := objects[1].(*widget.Label)
		borderContainer.Refresh()
		leftLabel.Text = fmt.Sprintf("%d x %s", tabItemWithAmount.amount, tabItemWithAmount.tabItem.Description)
		rightLabel.Text = fmt.Sprintf("%.2f", tabItemWithAmount.subTotal)
	})
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
