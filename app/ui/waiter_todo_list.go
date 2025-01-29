package ui

import (
	"fmt"
	"golangsevillabar/app/apiclient"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const WaiterTodoListStage = "WaiterTodoList"

type waiterTodoListScreen struct {
	container     *fyne.Container
	waiters       []string
	readApiClient *apiclient.ReadClient
	table         int
	tableLabel    *widget.Label
}

func (w *waiterTodoListScreen) ExecuteOnTakeOver(param interface{}) {
	tableNumber := param.(int)
	w.table = tableNumber
	w.tableLabel.Text = fmt.Sprintf("%d", w.table)
}

func (w *waiterTodoListScreen) GetPaintedContainer() *fyne.Container {
	return w.container
}

func (w *waiterTodoListScreen) GetStageName() string {
	return WaiterTodoListStage
}

func (w *waiterTodoListScreen) UpdateWaiterTodo() {
	response, err := w.readApiClient.GetTodoListForWaiter()

	if err != nil {
		slog.Error("error calling readApi", slog.Any("error", err))
	}

	if !response.OK {
		slog.Error("error calling serverApi", slog.Any("error", response.Error))
	}

	// todoItemsByTable := response.Data
}

func CreateWaiterTodoListScreen(waiters []string, readApiClient *apiclient.ReadClient) *waiterTodoListScreen {
	waiterDashContainer := container.NewVBox(widget.NewCard("Waiters", "", widget.NewForm()))

	return &waiterTodoListScreen{
		container:     waiterDashContainer,
		waiters:       waiters,
		readApiClient: readApiClient,
		tableLabel:    widget.NewLabel(""),
	}
}
