package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const InvoiceStage = "Invoice"

type invoiceScreen struct {
	table          int
	form           *widget.Form
	tableLabel     *widget.Label
	itemsList      *widget.List
	container      *fyne.Container
	readApiClient  *apiclient.ReadClient
	writeApiClient *apiclient.WriteClient
	stageManager   *StageManager
}

func (i *invoiceScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumber := param.(int)
	i.table = tableNumber

	response, err := i.readApiClient.GetInvoiceForTable(tableNumber)
	if err != nil {
		slog.Error("client error calling readapi", slog.Any("error", err))
		return
	}

	if !response.OK {
		slog.Error("server error calling readapi", slog.Any("error", err))
		return
	}

	invoice := response.Data

	for _, item := range invoice.TabInvoice.Items {
		i.itemsList.UpdateItem(1, widget.NewLabel(fmt.Sprintf("%d %s %f", item.MenuNumber, item.Description, item.Price)))
	}

	i.itemsList.Refresh()
	i.tableLabel.Text = fmt.Sprintf("%d", i.table)
}

func (i *invoiceScreen) GetPaintedContainer() *fyne.Container {
	return i.container
}

func (i *invoiceScreen) GetStageName() string {
	return InvoiceStage
}

func CreateInvoiceScreen(readApiClient *apiclient.ReadClient, writeApiClient *apiclient.WriteClient, stageManager *StageManager) *invoiceScreen {
	form := &widget.Form{}
	tableLabel := widget.NewLabel("")
	itemsList := widget.NewList(func() int { return 3 }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(lii widget.ListItemID, co fyne.CanvasObject) {})
	totalLabel := widget.NewLabel("")
	hasUnservedItemsLabel := widget.NewLabel("")
	form.Append("Table", tableLabel)
	form.Append("Items", itemsList)
	form.Append("Total", totalLabel)
	form.Append("Has Unserved Items?", hasUnservedItemsLabel)
	container := container.NewVBox()
	container.Add(widget.NewCard("Invoice", "", form))

	return &invoiceScreen{
		table:          0,
		form:           form,
		tableLabel:     tableLabel,
		container:      container,
		readApiClient:  readApiClient,
		writeApiClient: writeApiClient,
		stageManager:   stageManager,
		itemsList:      itemsList,
	}
}
