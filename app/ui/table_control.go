package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type tableControl struct {
	Container    *fyne.Container
	tableButtons []*tableButton
}

func CreateTableControl(totalTables int) *tableControl {

	tableButtons := []*tableButton{}
	grid := container.New(layout.NewGridLayout(3))
	for i := 0; i < totalTables; i++ {
		tableButton := newTableButton(i + 1)
		tableButtons = append(tableButtons, tableButton)
		grid.Add(tableButton)
	}

	return &tableControl{Container: grid, tableButtons: tableButtons}
}

func (tc *tableControl) UpdateActiveTables(activeTables []int) {
	for _, tb := range tc.tableButtons {
		tb.SetInactive()
	}
	for _, buttonID := range activeTables {
		tc.tableButtons[buttonID-1].SetActive()
	}
	tc.Container.Refresh()
}
