package ui

import (
	"golangsevillabar/app/apiclient"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type tableControl struct {
	Container    *fyne.Container
	tableButtons []*tableButton
	client       *apiclient.Client
	waiters      []string
}

func CreateTableControl(totalTables int, client *apiclient.Client, mainContainer *fyne.Container) *tableControl {

	tableButtons := []*tableButton{}
	grid := container.New(layout.NewGridLayout(3))
	for i := 0; i < totalTables; i++ {
		tableButton := newTableButton(i+1, []string{"w1", "w2"})
		tableButtons = append(tableButtons, tableButton)
		grid.Add(tableButton)
	}
	mainContainer.Add(grid)

	return &tableControl{Container: grid, tableButtons: tableButtons, client: client,
		waiters: []string{"w1", "w2"}}
}

func (tc *tableControl) UpdateActiveTables() {
	activeTables, err := tc.client.GetActiveTables()
	if err != nil {
		slog.Error("client error calling readapi", slog.Any("error", err))
		return
	}
	if !activeTables.OK {
		slog.Error("server error calling readapi", slog.Any("error", activeTables.Error))
		return
	}

	for _, tb := range tc.tableButtons {
		tb.SetInactive()
	}
	for _, buttonID := range activeTables.Data.ActiveTables {
		tc.tableButtons[buttonID-1].SetActive()
	}
	tc.Container.Refresh()
}
