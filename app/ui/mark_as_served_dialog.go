package ui

import (
	"cqrseventsourcingbar/app/apiclient"
	"cqrseventsourcingbar/queries"
	"cqrseventsourcingbar/writeservice/model"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createMarkAsServedDialog(tableNumberAndTabId tableNumberAndTabId, w fyne.Window, pendingTabItems []queries.TabItem, writeApiClient *apiclient.WriteClient, stageManager *StageManager) *dialog.FormDialog {
	metadataByOptions := make(map[string]tabItemWithAmount)

	pendingItemsWithAmount := getTabItemsWithAmount(pendingTabItems)

	var options []string

	for _, tabItemWithAmount := range pendingItemsWithAmount {
		option := fmt.Sprintf("%d x %s - menu number: %d", tabItemWithAmount.amount, tabItemWithAmount.tabItem.Description, tabItemWithAmount.tabItem.MenuNumber)
		options = append(options, option)
		metadataByOptions[option] = tabItemWithAmount
	}

	group := widget.NewCheckGroup(options, func(s []string) {})
	group.Required = true
	group.SetSelected(options)

	item := widget.NewFormItem("", group)
	formItems := []*widget.FormItem{item}

	return dialog.NewForm(fmt.Sprintf("Mark as served items for table %d", tableNumberAndTabId.tableNumber), "Confirm", "Cancel", formItems, func(confirm bool) {
		if confirm {
			checkGroup := formItems[0].Widget.(*widget.CheckGroup)
			selectedOptions := checkGroup.Selected
			selectedMenuNumbers := []int{}

			for _, selectedOption := range selectedOptions {
				metadata := metadataByOptions[selectedOption]
				for i := 0; i < metadata.amount; i++ {
					selectedMenuNumbers = append(selectedMenuNumbers, metadata.tabItem.MenuNumber)
				}
			}

			err := writeApiClient.ExecuteCommand(model.MarkDrinksServedRequest{
				TabId:       tableNumberAndTabId.tabId,
				MenuNumbers: selectedMenuNumbers,
			})

			if err != nil {
				slog.Error("error calling writeApi with MarkDrinksServedRequest", slog.Any("error", err))
				return
			}

			err = stageManager.TakeOver(MainContentStage, nil)
			if err != nil {
				slog.Error("error launching main content screen", slog.Any("error", err))
			}
		}
	}, w)
}
