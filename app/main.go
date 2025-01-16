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

	apiClient := apiclient.NewClient(&http.Client{}, "http://localhost:8081")

	// tabIdForTableOne, err := apiClient.GetTabIdForTable(1)
	// if err != nil {
	// 	slog.Error("error from api", slog.Any("error", err))
	// }
	// slog.Info("TabId for table 1", slog.Any("ID", tabIdForTableOne.Data))

	stageManager := ui.CreateStageManager()

	mainContainer := ui.CreateMainContent(amountOfTables, apiClient, &stageManager, waiters, w)
	openTabStage := ui.CreateOpenTabScreen(waiters, &stageManager)

	stageManager.RegisterStager(mainContainer)
	stageManager.RegisterStager(openTabStage)

	// tableControl := ui.CreateTableControl(6, apiClient, w)

	// tableControl.UpdateActiveTables()

	w.SetContent(stageManager.GetContainer())

	err := stageManager.TakeOver(ui.MainContentStage, nil)
	if err != nil {
		slog.Error("error opening screen", slog.Any("error", err))
	}

	// err = stageManager.TakeOver(ui.OpenTabStage, 4)
	// if err != nil {
	// 	slog.Error("error opening screen", slog.Any("error", err))
	// }

	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
