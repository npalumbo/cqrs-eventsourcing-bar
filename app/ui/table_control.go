package ui

import (
	"golangsevillabar/app/apiclient"
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myTheme struct{}

// Color implements fyne.Theme.
func (m *myTheme) Color(color fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(color, variant)
}

// Icon implements fyne.Theme.
func (m *myTheme) Icon(iconName fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(iconName)
}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInputRadius {
		return 100
	}
	return theme.DefaultTheme().Size(name)
}

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
		grid.Add(container.NewThemeOverride(tableButton, &myTheme{}))
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
