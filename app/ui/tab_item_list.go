package ui

import (
	"fmt"

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
