package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/queries"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const InvoiceStage = "Invoice"

type invoiceScreen struct {
	table                 int
	containerInCard       *fyne.Container
	tableLabel            *widget.Label
	totalLabel            *widget.Label
	hasUnservedItemsLabel *widget.Label
	itemsList             *widget.List
	container             *fyne.Container
	readApiClient         *apiclient.ReadClient
	writeApiClient        *apiclient.WriteClient
	stageManager          *StageManager
	tabItems              *[]queries.TabItem
	closeTabButton        *widget.Button
	closeTabForm          *dialog.FormDialog
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

	*i.tabItems = nil
	*i.tabItems = append(*i.tabItems, invoice.TabInvoice.Items...)

	if invoice.TabInvoice.HasUnservedItems {
		i.closeTabButton.Disable()
	}

	i.itemsList.Refresh()
	i.containerInCard.Refresh()
	i.container.Refresh()
	i.tableLabel.Text = fmt.Sprintf("%d", i.table)
	i.totalLabel.Text = fmt.Sprintf("%.2f", invoice.TabInvoice.Total)
	i.hasUnservedItemsLabel.Text = fmt.Sprintf("%t", invoice.TabInvoice.HasUnservedItems)
}

func (i *invoiceScreen) GetPaintedContainer() *fyne.Container {
	return i.container
}

func (i *invoiceScreen) GetStageName() string {
	return InvoiceStage
}

func CreateInvoiceScreen(readApiClient *apiclient.ReadClient, writeApiClient *apiclient.WriteClient, stageManager *StageManager, w fyne.Window) *invoiceScreen {
	tableLabel := widget.NewLabel("")
	tabItems := &[]queries.TabItem{}
	itemsList := widget.NewList(func() int { return len(*tabItems) }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(lii widget.ListItemID, co fyne.CanvasObject) {
		tabItem := (*tabItems)[lii]
		listItemLabel := co.(*widget.Label)
		listItemLabel.Text = fmt.Sprintf("%s %.2f", tabItem.Description, tabItem.Price)
		listItemLabel.Refresh()
	})
	totalLabel := widget.NewLabel("")
	hasUnservedItemsLabel := widget.NewLabel("")
	formItems := []*widget.FormItem{}
	formItems = append(formItems, widget.NewFormItem("Total", totalLabel))
	formItems = append(formItems, widget.NewFormItem("Tip", widget.NewEntry()))
	closeTabForm := dialog.NewForm("Close Tab", "Close", "Cancel", formItems, func(b bool) {}, w)
	closeTabButton := widget.NewButton("Close Tab", func() {
		closeTabForm.Show()
	})
	containerInCard := container.NewGridWithColumns(2,
		widget.NewLabel("Table"), tableLabel,
		widget.NewLabel("Items"), itemsList,
		widget.NewLabel("Total"), totalLabel,
		widget.NewLabel("Has UnservedItems"), hasUnservedItemsLabel,
		widget.NewButton("Back", func() {
			err := stageManager.TakeOver(MainContentStage, nil)
			if err != nil {
				slog.Error("error launching main content screen", slog.Any("error", err))
			}
		}),
		closeTabButton,
	)

	container := container.NewStack(widget.NewCard("Invoice", "", containerInCard))

	return &invoiceScreen{
		table:                 0,
		containerInCard:       containerInCard,
		tableLabel:            tableLabel,
		totalLabel:            totalLabel,
		hasUnservedItemsLabel: hasUnservedItemsLabel,
		container:             container,
		readApiClient:         readApiClient,
		writeApiClient:        writeApiClient,
		stageManager:          stageManager,
		itemsList:             itemsList,
		tabItems:              tabItems,
		closeTabButton:        closeTabButton,
		closeTabForm:          closeTabForm,
	}
}
