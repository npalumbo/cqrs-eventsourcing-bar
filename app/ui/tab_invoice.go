package ui

import (
	"errors"
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/queries"
	"golangsevillabar/writeservice/model"
	"log/slog"
	"strconv"

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
	tipLabel              *widget.Label
	itemsList             *widget.List
	container             *fyne.Container
	readApiClient         *apiclient.ReadClient
	writeApiClient        *apiclient.WriteClient
	stageManager          *StageManager
	tabItems              *[]queries.TabItem
	closeTabButton        *widget.Button
	closeTabForm          *dialog.FormDialog
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
	// i.payingWithEntry.SetValidationError(errors.New("must set paying with"))

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

	*i.tabItems = nil
	*i.tabItems = append(*i.tabItems, invoice.Items...)

	if invoice.HasUnservedItems {
		i.closeTabButton.Disable()
	}

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
	payingWithEntry := widget.NewEntry()

	payingWithFormItem := widget.NewFormItem("Paying with", payingWithEntry)
	payingWithFormItem.HintText = "Amount"

	formItems = append(formItems, payingWithFormItem)

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
		tabItems:              tabItems,
		payingWithEntry:       payingWithEntry,
		currentTotal:          -1,
		currentTip:            0,
	}

	closeTabForm := dialog.NewForm("Close Tab", "Close", "Cancel", formItems, func(hitCloseButton bool) {
		slog.Info("Value of b", slog.Any("b", hitCloseButton))

		if hitCloseButton {

			amount, err := strconv.ParseFloat(invoiceScreen.payingWithEntry.Text, 64)
			if err != nil {
				slog.Error("error converting paying with before closing tab", slog.Any("error", err))
			}
			err = writeApiClient.ExecuteCommand(model.CloseTabRequest{
				TabId:      invoiceScreen.currentInvoiceData.TabID,
				AmountPaid: amount,
			})
			if err != nil {
				slog.Error("error calling write api", slog.Any("error", err))
			}
			// If no error, we asume the close tab command worked and refresh the Tip field
			invoiceScreen.currentTip = amount - invoiceScreen.currentTotal
			invoiceScreen.tipLabel.Text = fmt.Sprintf("%.2f", invoiceScreen.currentTip)
			invoiceScreen.closeTabButton.Disable()
			invoiceScreen.tipLabel.Refresh()
			invoiceScreen.containerInCard.Refresh()
		}
	}, w)

	closeTabForm.Refresh()

	closeTabButton := widget.NewButton("Close Tab", func() {
		closeTabForm.Show()
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

	payingWithEntry.SetValidationError(errors.New("must set paying with"))

	payingWithEntry.Validator = func(s string) error {
		amount, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		if amount < invoiceScreen.currentTotal {
			return errors.New("need to pay with an amount higher than total")
		}

		return nil
	}

	invoiceScreen.containerInCard = containerInCard
	invoiceScreen.container = container
	invoiceScreen.closeTabButton = closeTabButton
	invoiceScreen.closeTabForm = closeTabForm

	return invoiceScreen
}
