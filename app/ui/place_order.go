package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/shared"
	"golangsevillabar/writeservice/model"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const PlaceOrderStage = "PlaceOrder"

type placeOrderScreen struct {
	container            *fyne.Container
	placeOrderScreenCard *widget.Card
	table                int
	writeApiClient       *apiclient.WriteClient
	readApiClient        *apiclient.ReadClient
	stageManager         *StageManager
	form                 *widget.Form
	allMenuItems         []shared.MenuItem
	tabId                *string
}

func (p *placeOrderScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumberAndTabId := param.(tableNumberAndTabId)

	for _, formItem := range p.form.Items {
		selectAmount := formItem.Widget.(*widget.Select)
		selectAmount.SetSelected("0")
	}
	p.table = tableNumberAndTabId.tableNumber
	p.placeOrderScreenCard.Title = fmt.Sprintf("Order drinks for table %d", p.table)
	p.tabId = &tableNumberAndTabId.tabId
}

func (p *placeOrderScreen) GetPaintedContainer() *fyne.Container {
	return p.container
}

func (p *placeOrderScreen) GetStageName() string {
	return PlaceOrderStage
}

func CreatePlaceOrderScreen(writeApiClient *apiclient.WriteClient, readApiClient *apiclient.ReadClient, stageManager *StageManager) *placeOrderScreen {
	container := container.NewStack()

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
	placeOrderScreenCard := widget.NewCard("Order drinks", "", form)

	placeOrderScreen := &placeOrderScreen{
		container:            container,
		placeOrderScreenCard: placeOrderScreenCard,
		writeApiClient:       writeApiClient,
		readApiClient:        readApiClient,
		stageManager:         stageManager,
		form:                 form,
		allMenuItems:         allMenuItems,
	}

	form.SubmitText = "OK"
	form.CancelText = "Cancel"
	form.OnCancel = func() {
		err := stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			slog.Error("error opening main content screen", slog.Any("error", err))
		}
	}

	form.OnSubmit = func() {

		orderedItems := []int{}
		for i, formItem := range placeOrderScreen.form.Items {
			selectAmount := formItem.Widget.(*widget.Select)

			amount, err := strconv.Atoi(selectAmount.Selected)

			if err != nil {
				slog.Error("error converting amount selected of menuItem", slog.Any("error", err), slog.Any("menu_item", i+1))
				return
			}

			for j := 0; j < amount; j++ {
				orderedItems = append(orderedItems, i+1)
			}
		}

		err = writeApiClient.ExecuteCommand(model.PlaceOrderRequest{
			TabId:     *placeOrderScreen.tabId,
			MenuItems: orderedItems,
		})

		if err != nil {
			slog.Error("client error calling writeapi", slog.Any("error", err))
			return
		}

		err := stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			slog.Error("error opening main content screen", slog.Any("error", err))
		}
	}

	container.Add(placeOrderScreenCard)
	return placeOrderScreen
}
