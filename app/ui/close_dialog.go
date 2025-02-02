package ui

import (
	"cqrseventsourcingbar/app/apiclient"
	"cqrseventsourcingbar/writeservice/model"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createCoseTabFormDialog(w fyne.Window, payingWithEntry *widget.Entry, invoiceScreen *invoiceScreen, writeApiClient *apiclient.WriteClient) *dialog.FormDialog {
	formItems := []*widget.FormItem{}
	formItems = append(formItems, widget.NewFormItem("Total", widget.NewLabel(fmt.Sprintf("%.2f", invoiceScreen.currentTotal))))

	payingWithFormItem := widget.NewFormItem("Paying with", payingWithEntry)
	payingWithFormItem.HintText = "Amount"

	formItems = append(formItems, payingWithFormItem)
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
	closeTabDialog := dialog.NewForm("Close Tab", "Close", "Cancel", formItems, func(hitCloseButton bool) {
		if hitCloseButton {

			amount, err := strconv.ParseFloat(payingWithEntry.Text, 64)
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

	return closeTabDialog
}
