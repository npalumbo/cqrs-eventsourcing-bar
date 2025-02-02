package main

import (
	"golangsevillabar/app/apiclient"
	"golangsevillabar/app/ui"
	"log/slog"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("CQRS ES BAR")

	const amountOfTables = 6
	waiters := []string{"w1", "w2"}

	readApiClient := apiclient.NewReadClient(&http.Client{}, "http://localhost:8081")
	writeApiClient := apiclient.NewWriteClient(&http.Client{}, "http://localhost:8080")

	stageManager := ui.CreateStageManager()

	tableControl := ui.CreateTableControl(amountOfTables, readApiClient, &stageManager)
	waiterControl := ui.CreateWaiterControl(readApiClient, writeApiClient, waiters, &w, &stageManager)
	mainContainer := ui.CreateMainContentScreen(tableControl, waiterControl)
	openTabStage := ui.CreateOpenTabScreen(waiters, writeApiClient, &stageManager)
	invoiceStage := ui.CreateInvoiceScreen(readApiClient, writeApiClient, &stageManager, w)
	placeOrderStage := ui.CreatePlaceOrderScreen(writeApiClient, readApiClient, &stageManager)
	tabStatusStage := ui.CreateTabStatusScreen(&stageManager)

	stageManager.RegisterStager(mainContainer)
	stageManager.RegisterStager(openTabStage)
	stageManager.RegisterStager(invoiceStage)
	stageManager.RegisterStager(placeOrderStage)
	stageManager.RegisterStager(tabStatusStage)

	w.SetContent(stageManager.GetContainer())

	err := stageManager.TakeOver(ui.MainContentStage, nil)
	if err != nil {
		slog.Error("error opening main content screen", slog.Any("error", err))
	}

	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
