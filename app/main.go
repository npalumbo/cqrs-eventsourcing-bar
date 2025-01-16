package main

import (
	"golangsevillabar/app/apiclient"
	"golangsevillabar/app/ui"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("CQRS ES BAR")

	apiClient := apiclient.NewClient(&http.Client{}, "http://localhost:8081")

	// tabIdForTableOne, err := apiClient.GetTabIdForTable(1)
	// if err != nil {
	// 	slog.Error("error from api", slog.Any("error", err))
	// }
	// slog.Info("TabId for table 1", slog.Any("ID", tabIdForTableOne.Data))

	mainContainer := ui.CreateMainContent(6, apiClient, w)

	// tableControl := ui.CreateTableControl(6, apiClient, w)

	// tableControl.UpdateActiveTables()

	w.SetContent(mainContainer.MakeUI())

	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
