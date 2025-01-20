package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const InvoiceStage = "Invoice"

type invoiceScreen struct {
	table          int
	form           *widget.Form
	tableLabel     *widget.Label
	container      *fyne.Container
	writeApiClient *apiclient.WriteClient
	stageManager   *StageManager
}

func (i *invoiceScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumber := param.(int)
	i.table = tableNumber
	i.tableLabel.Text = fmt.Sprintf("%d", i.table)
}

func (i *invoiceScreen) GetPaintedContainer() *fyne.Container {
	return i.container
}

func (i *invoiceScreen) GetStageName() string {
	return InvoiceStage
}

func CreateInvoiceScreen(writeApiClient *apiclient.WriteClient, stageManager *StageManager) *invoiceScreen {
	form := &widget.Form{}
	tableLabel := widget.NewLabel("")
	itemsLabel := widget.NewLabel("")
	totalLabel := widget.NewLabel("")
	hasUnservedItemsLabel := widget.NewLabel("")
	form.Append("Table", tableLabel)
	form.Append("Items", itemsLabel)
	form.Append("Total", totalLabel)
	form.Append("Has Unserved Items?", hasUnservedItemsLabel)
	container := container.NewVBox()
	container.Add(widget.NewCard("Invoice", "", form))

	return &invoiceScreen{
		table:          0,
		form:           form,
		tableLabel:     tableLabel,
		container:      container,
		writeApiClient: writeApiClient,
		stageManager:   stageManager,
	}
}
