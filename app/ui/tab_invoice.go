package ui

import (
	"cqrseventsourcingbar/app/apiclient"
	"cqrseventsourcingbar/queries"
	"fmt"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const InvoiceStage = "Invoice"

type invoiceScreen struct {
	table                 int
	containerInCard       *fyne.Container
	invoiceScreenCard     *widget.Card
	totalLabel            *widget.Label
	hasUnservedItemsLabel *widget.Label
	tipLabel              *widget.Label
	itemsList             *widget.List
	container             *fyne.Container
	readApiClient         *apiclient.ReadClient
	writeApiClient        *apiclient.WriteClient
	stageManager          *StageManager
	tabItemsWithAmount    *[]tabItemWithAmount
	closeTabButton        *widget.Button
	payingWithEntry       *widget.Entry
	currentTotal          float64
	currentTip            float64
	currentInvoiceData    *queries.TabInvoice
}

func (i *invoiceScreen) ExecuteOnTakeOver(param interface{}) {
	i.currentTotal = -1
	i.currentTip = 0
	i.tipLabel.Text = ""
	tableNumber := param.(int)
	i.table = tableNumber
	i.payingWithEntry.Text = "0"

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
	i.currentInvoiceData = &invoice

	*i.tabItemsWithAmount = nil
	*i.tabItemsWithAmount = append(*i.tabItemsWithAmount, getTabItemsWithAmount(invoice.Items)...)

	i.itemsList.Refresh()
	i.containerInCard.Refresh()
	i.container.Refresh()
	i.invoiceScreenCard.SetTitle(fmt.Sprintf("Invoice for table %d", i.table))
	i.totalLabel.Text = fmt.Sprintf("%.2f", invoice.Total)
	i.currentTotal, err = strconv.ParseFloat(fmt.Sprintf("%.2f", invoice.Total), 64)
	if err != nil {
		slog.Error("could not convert current total to float", slog.Any("error", err))
	}

	i.hasUnservedItemsLabel.Text = fmt.Sprintf("%t", invoice.HasUnservedItems)
	if invoice.HasUnservedItems {
		i.closeTabButton.Disable()
	} else {
		i.closeTabButton.Enable()
	}
}

func (i *invoiceScreen) GetPaintedContainer() *fyne.Container {
	return i.container
}

func (i *invoiceScreen) GetStageName() string {
	return InvoiceStage
}

func CreateInvoiceScreen(readApiClient *apiclient.ReadClient, writeApiClient *apiclient.WriteClient, stageManager *StageManager, w fyne.Window) *invoiceScreen {
	tipLabel := widget.NewLabel("")
	tabItemsWithAmount := &[]tabItemWithAmount{}
	itemsList := CreateTabItemList(tabItemsWithAmount)
	totalLabel := widget.NewLabel("")
	hasUnservedItemsLabel := widget.NewLabel("")
	payingWithEntry := widget.NewEntry()
	payingWithEntry.Text = "0"

	invoiceScreen := &invoiceScreen{
		table:                 0,
		totalLabel:            totalLabel,
		tipLabel:              tipLabel,
		hasUnservedItemsLabel: hasUnservedItemsLabel,
		readApiClient:         readApiClient,
		writeApiClient:        writeApiClient,
		stageManager:          stageManager,
		itemsList:             itemsList,
		tabItemsWithAmount:    tabItemsWithAmount,
		payingWithEntry:       payingWithEntry,
		currentTotal:          -1,
		currentTip:            0,
	}

	closeTabButton := widget.NewButton("Close Tab", func() {
		closeTabDialog := createCoseTabFormDialog(w, payingWithEntry, invoiceScreen, writeApiClient)
		closeTabDialog.Show()
	})

	containerInCard := container.NewBorder(nil, container.NewGridWithRows(1,
		widget.NewButton("Back", func() {
			err := stageManager.TakeOver(MainContentStage, nil)
			if err != nil {
				slog.Error("error launching main content screen", slog.Any("error", err))
			}
		}),
		closeTabButton),
		nil, nil,
		container.NewGridWithColumns(2,
			widget.NewLabel("Items"), itemsList,
			widget.NewLabel("Total"), totalLabel,
			widget.NewLabel("Has UnservedItems"), hasUnservedItemsLabel,
			widget.NewLabel("Tip"), tipLabel,
		))

	invoiceScreenCard := widget.NewCard("Invoice", "", containerInCard)
	container := container.NewStack(invoiceScreenCard)

	invoiceScreen.containerInCard = containerInCard
	invoiceScreen.invoiceScreenCard = invoiceScreenCard
	invoiceScreen.container = container
	invoiceScreen.closeTabButton = closeTabButton

	return invoiceScreen
}
