package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/shared"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const PlaceOrderStage = "PlaceOrder"

type placeOrderScreen struct {
	container      *fyne.Container
	table          int
	tableLabel     *widget.Label
	writeApiClient *apiclient.WriteClient
	readApiClient  *apiclient.ReadClient
	stageManager   *StageManager
	form           *widget.Form
	allMenuItems   []shared.MenuItem
}

func (p *placeOrderScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumber := param.(int)
	p.table = tableNumber
	p.tableLabel.Text = fmt.Sprintf("%d", p.table)

	// tabId, err := p.readApiClient.GetTabIdForTable(tableNumber)

	// if err != nil {
	// 	slog.Error("client error calling readapi", slog.Any("error", err))
	// 	return
	// }

}

func (p *placeOrderScreen) GetPaintedContainer() *fyne.Container {
	return p.container
}

func (p *placeOrderScreen) GetStageName() string {
	return PlaceOrderStage
}

func CreatePlaceOrderScreen(writeApiClient *apiclient.WriteClient, readApiClient *apiclient.ReadClient, stageManager *StageManager) *placeOrderScreen {
	container := container.NewVBox()

	allMenuItemsResponse, err := readApiClient.GetAllMenuItems()

	if err != nil {
		slog.Error("client error calling readapi", slog.Any("error", err))
		return nil
	}

	menuFormItems := []*widget.FormItem{}

	allMenuItems := allMenuItemsResponse.Data

	for _, menuItem := range allMenuItems {
		selectWidget := widget.NewSelect([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}, func(s string) {})
		selectWidget.SetSelected("0")
		formItem := widget.NewFormItem(fmt.Sprintf("%s: %.2f", menuItem.Description, menuItem.Price), selectWidget)
		menuFormItems = append(menuFormItems, formItem)
	}

	form := widget.NewForm(menuFormItems...)
	form.SubmitText = "OK"
	form.CancelText = "Cancel"
	form.OnCancel = func() {
		err := stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			if err != nil {
				slog.Error("error opening main content screen", slog.Any("error", err))
			}
		}
	}

	container.Add(widget.NewCard("Order drinks", "", form))

	return &placeOrderScreen{
		container:      container,
		tableLabel:     &widget.Label{},
		writeApiClient: writeApiClient,
		readApiClient:  readApiClient,
		stageManager:   stageManager,
		form:           form,
		allMenuItems:   allMenuItems,
	}
}
