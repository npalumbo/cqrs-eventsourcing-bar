package ui

import (
	"cqrseventsourcingbar/app/apiclient"
	"cqrseventsourcingbar/writeservice/model"
	"fmt"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const OpenTabStage = "OpenTab"

type openTabScreen struct {
	table          int
	form           *widget.Form
	tableLabel     *widget.Label
	container      *fyne.Container
	writeApiClient *apiclient.WriteClient
	stageManager   *StageManager
}

func (o *openTabScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumber := param.(int)
	o.table = tableNumber
	o.tableLabel.Text = fmt.Sprintf("%d", o.table)
}

func (o *openTabScreen) GetPaintedContainer() *fyne.Container {
	return o.container
}

func (o *openTabScreen) GetStageName() string {
	return OpenTabStage
}

func CreateOpenTabScreen(waiters []string, writeApiClient *apiclient.WriteClient, stageManager *StageManager) *openTabScreen {
	form := &widget.Form{}
	tableLabel := widget.NewLabel("")
	waitersDropDown := widget.NewSelect(waiters, func(s string) {})
	waitersDropDown.SetSelected(waiters[0])
	form.Append("Table", tableLabel)
	form.Append("waiter", waitersDropDown)
	form.CancelText = "Back"
	form.OnCancel = func() {
		err := stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			slog.Error("error launching main content screen", slog.Any("error", err))
		}
	}
	form.SubmitText = "Open tab"
	form.OnSubmit = func() {
		tableNumber, err := strconv.Atoi(tableLabel.Text)
		if err != nil {
			slog.Error("error getting tableNumber", slog.Any("error", err))
		}
		err = writeApiClient.ExecuteCommand(model.OpenTabRequest{
			TableNumber: tableNumber,
			Waiter:      waitersDropDown.Selected,
		})
		if err != nil {
			slog.Error("error sending command", slog.Any("error", err))
		}
		err = stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			slog.Error("error launching main content screen", slog.Any("error", err))
		}

	}

	container := container.NewVBox()
	container.Add(widget.NewCard("Open a Tab", "", form))

	return &openTabScreen{
		form:           form,
		tableLabel:     tableLabel,
		container:      container,
		writeApiClient: writeApiClient,
		stageManager:   stageManager,
	}
}
