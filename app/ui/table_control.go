package ui

import (
	"cqrseventsourcingbar/app/apiclient"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type tableControl struct {
	Card           *widget.Card
	innerContainer *fyne.Container
	tableButtons   []*tableButton
	client         *apiclient.ReadClient
}

func CreateTableControl(totalTables int, client *apiclient.ReadClient, stageManager *StageManager) *tableControl {

	tableButtons := []*tableButton{}
	grid := container.New(layout.NewGridLayout(3))
	for i := 0; i < totalTables; i++ {
		tableButton := newTableButton(i+1, stageManager)
		tableButtons = append(tableButtons, tableButton)
		grid.Add(container.NewThemeOverride(tableButton, &roundButtonTheme{}))
	}
	card := widget.NewCard("Table Control", "", grid)

	return &tableControl{Card: card, innerContainer: grid, tableButtons: tableButtons, client: client}
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
		tc.tableButtons[tableID-1].SetActive(&tabStatus.Data)
	}
	tc.innerContainer.Refresh()
}
