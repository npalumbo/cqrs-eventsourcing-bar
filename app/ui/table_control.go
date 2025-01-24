package ui

import (
	"golangsevillabar/app/apiclient"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type tableControl struct {
	Container    *fyne.Container
	tableButtons []*tableButton
	client       *apiclient.ReadClient
	waiters      []string
}

func CreateTableControl(totalTables int, client *apiclient.ReadClient, waiters []string, stageManager *StageManager, mainContainer *fyne.Container) *tableControl {

	tableButtons := []*tableButton{}
	grid := container.New(layout.NewGridLayout(3))
	for i := 0; i < totalTables; i++ {
		tableButton := newTableButton(i+1, waiters, stageManager)
		tableButtons = append(tableButtons, tableButton)
		grid.Add(tableButton)
	}
	mainContainer.Add(widget.NewCard("Table Control", "", grid))

	return &tableControl{Container: grid, tableButtons: tableButtons, client: client,
		waiters: waiters}
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
	for _, tableID := range activeTables.Data {
		tabStatus, err := tc.client.GetTabForTable(tableID)
		if err != nil {
			slog.Error("client error calling readapi", slog.Any("error", err))
			return
		}
		if !tabStatus.OK {
			slog.Error("server error calling readapi", slog.Any("error", activeTables.Error))
			return
		}
		tc.tableButtons[tableID-1].SetActive(&tabStatus)
	}
	tc.Container.Refresh()
}
