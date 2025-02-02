package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"golangsevillabar/queries"
	"log/slog"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type waiterControl struct {
	Card               *widget.Card
	container          *fyne.Container
	waiters            []string
	client             *apiclient.ReadClient
	writeClient        *apiclient.WriteClient
	containersByWaiter map[string]*fyne.Container
	window             *fyne.Window
	stageManager       *StageManager
}

func (wc *waiterControl) UpdateWaiterControl() {
	for _, waiter := range wc.waiters {

		response, err := wc.client.GetTodoListForWaiter(waiter)

		if err != nil {
			slog.Error(fmt.Sprintf("client error when calling GetTodoListForWaiter, for waiter: %s", waiter), slog.Any("error", err))
			return
		}

		if !response.OK {
			slog.Error(fmt.Sprintf("server error calling GetTodoListForWaiter, for waiter: %s", waiter), slog.Any("error", response.Error))
			return
		}

		tabItemsByTable := response.Data

		container := wc.containersByWaiter[waiter]

		container.Objects = nil
		container.Add(widget.NewLabel(waiter))
		sortedKeys := getSortedKeys(tabItemsByTable)

		for _, table := range sortedKeys {
			itemsForTable := tabItemsByTable[table]
			tableButton := widget.NewButton(fmt.Sprintf("Table %d", table), nil)
			if len(itemsForTable) == 0 {
				tableButton.Disable()
			} else {
				tabIdForTableResponse, err := wc.client.GetTabIdForTable(table)

				if err != nil {
					slog.Error("client error calling readapi", slog.Any("error", err))
					continue
				}

				if !tabIdForTableResponse.OK {
					slog.Error("server error calling readapi", slog.Any("error", err))
					continue
				}

				tableButton.OnTapped = func() {
					dialog := createMarkAsServedDialog(tableNumberAndTabId{
						tableNumber: table,
						tabId:       tabIdForTableResponse.Data,
					}, *wc.window, itemsForTable, wc.writeClient, wc.stageManager)
					dialog.Show()

				}
			}
			container.Add(tableButton)
		}
	}
}

func getSortedKeys(tabItemsByTable map[int][]queries.TabItem) []int {
	keys := make([]int, 0, len(tabItemsByTable))
	for k := range tabItemsByTable {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func CreateWaiterControl(client *apiclient.ReadClient, writeClient *apiclient.WriteClient, waiters []string, w *fyne.Window, stageManager *StageManager) *waiterControl {

	containersByWaiter := make(map[string]*fyne.Container)

	grid := container.NewGridWithRows(len(waiters))
	for _, waiter := range waiters {
		containerForWaiter := container.NewHBox()
		containersByWaiter[waiter] = containerForWaiter
		grid.Add(containerForWaiter)
	}

	card := widget.NewCard("Waiter Control", "", grid)

	return &waiterControl{
		Card:               card,
		container:          grid,
		waiters:            waiters,
		client:             client,
		writeClient:        writeClient,
		containersByWaiter: containersByWaiter,
		window:             w,
		stageManager:       stageManager,
	}

}
