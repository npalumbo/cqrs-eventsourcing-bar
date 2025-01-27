package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/queries"
	"log/slog"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const InvoiceStage = "Invoice"

type tabItemWithAmount struct {
	tabItem  queries.TabItem
	amount   int
	subTotal float64
}

type invoiceScreen struct {
	table                 int
	containerInCard       *fyne.Container
	tableLabel            *widget.Label
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
	i.payingWithEntry.Text = ""

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
	*i.tabItemsWithAmount = append(*i.tabItemsWithAmount, tabItemsWithAmount(invoice.Items)...)

	i.itemsList.Refresh()
	i.containerInCard.Refresh()
	i.container.Refresh()
	i.tableLabel.Text = fmt.Sprintf("%d", i.table)
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
	tableLabel := widget.NewLabel("")
	tipLabel := widget.NewLabel("")
	tabItemsWithAmount := &[]tabItemWithAmount{}
	itemsList := widget.NewList(func() int { return len(*tabItemsWithAmount) }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(lii widget.ListItemID, co fyne.CanvasObject) {
		tabItemWithAmount := (*tabItemsWithAmount)[lii]
		listItemLabel := co.(*widget.Label)
		listItemLabel.Text = fmt.Sprintf("%d x %s %.2f  ...  %.2f", tabItemWithAmount.amount, tabItemWithAmount.tabItem.Description, tabItemWithAmount.tabItem.Price, tabItemWithAmount.subTotal)
		listItemLabel.Refresh()
	})
	totalLabel := widget.NewLabel("")
	hasUnservedItemsLabel := widget.NewLabel("")
	payingWithEntry := widget.NewEntry()

	invoiceScreen := &invoiceScreen{
		table:                 0,
		tableLabel:            tableLabel,
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

	closeTabDialog := createCoseTabFormDialog(w, payingWithEntry, invoiceScreen, writeApiClient)

	closeTabDialog.Refresh()

	closeTabButton := widget.NewButton("Close Tab", func() {
		closeTabDialog.Show()
	})

	containerInCard := container.NewGridWithColumns(2,
		widget.NewLabel("Table"), tableLabel,
		widget.NewLabel("Items"), itemsList,
		widget.NewLabel("Total"), totalLabel,
		widget.NewLabel("Has UnservedItems"), hasUnservedItemsLabel,
		widget.NewLabel("Tip"), tipLabel,
		widget.NewButton("Back", func() {
			err := stageManager.TakeOver(MainContentStage, nil)
			if err != nil {
				slog.Error("error launching main content screen", slog.Any("error", err))
			}
		}),
		closeTabButton,
	)

	container := container.NewStack(widget.NewCard("Invoice", "", containerInCard))

	invoiceScreen.containerInCard = containerInCard
	invoiceScreen.container = container
	invoiceScreen.closeTabButton = closeTabButton

	return invoiceScreen
}

func tabItemsWithAmount(tabItems []queries.TabItem) []tabItemWithAmount {
	var amountsPerItemID map[int]int = make(map[int]int)
	tabItemMenuNumbers := []int{}
	var tabItemsByMenuNumbers map[int]queries.TabItem = make(map[int]queries.TabItem)

	for _, tabItem := range tabItems {
		_, ok := tabItemsByMenuNumbers[tabItem.MenuNumber]
		if !ok {
			tabItemsByMenuNumbers[tabItem.MenuNumber] = tabItem
		}
		amount, ok := amountsPerItemID[tabItem.MenuNumber]
		if !ok {
			amountsPerItemID[tabItem.MenuNumber] = 1
		} else {
			amountsPerItemID[tabItem.MenuNumber] = amount + 1
		}
		tabItemMenuNumbers = append(tabItemMenuNumbers, tabItem.MenuNumber)
	}

	slices.Sort(tabItemMenuNumbers)
	tabItemWithAmounts := []tabItemWithAmount{}

	for _, menuNumber := range tabItemMenuNumbers {
		amount := amountsPerItemID[menuNumber]
		tabItemWithAmounts = append(tabItemWithAmounts, tabItemWithAmount{
			amount:   amount,
			tabItem:  tabItemsByMenuNumbers[menuNumber],
			subTotal: float64(amount) * tabItemsByMenuNumbers[menuNumber].Price,
		})
	}

	return tabItemWithAmounts
}
