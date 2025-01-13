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

	apiClient := apiclient.NewClient(&http.Client{}, "http://localhost:8081")

	active, err := apiClient.GetActiveTables()
	if err != nil {
		slog.Error("error from api", slog.Any("error", err))
	}

	tabIdForTableOne, err := apiClient.GetTabIdForTable(1)
	if err != nil {
		slog.Error("error from api", slog.Any("error", err))
	}
	slog.Info("TabId for table 1", slog.Any("ID", tabIdForTableOne.Data))

	buttonControl := ui.CreateTableControl(6)

	w.SetContent(buttonControl.Container)
	buttonControl.UpdateActiveTables(active.Data.ActiveTables)

	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
