package ui

import (
	"golangsevillabar/app/apiclient"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const MainContentStage = "MainContent"

type MainContent struct {
	tableControl         *tableControl
	waiterControl        *waiterControl
	stagerManager        *StageManager
	mainContentContainer *fyne.Container
}

func (m *MainContent) ExecuteOnTakeOver(param interface{}) {
	time.Sleep(200 * time.Millisecond)
	m.tableControl.UpdateActiveTables()
	m.waiterControl.UpdateWaiterControl()
}

func (m *MainContent) GetPaintedContainer() *fyne.Container {
	return m.mainContentContainer
}

func (m *MainContent) GetStageName() string {
	return MainContentStage
}

func CreateMainContent(totalTables int, client *apiclient.ReadClient, writeClient *apiclient.WriteClient, stageManager *StageManager, waiters []string, w fyne.Window) *MainContent {

	tableControl := CreateTableControl(totalTables, client, stageManager)
	waiterControl := CreateWaiterControl(client, writeClient, waiters, &w, stageManager)

	mainContentContainer := container.NewBorder(nil, waiterControl.Card, nil, nil, tableControl.Card)

	tableControl.UpdateActiveTables()
	waiterControl.UpdateWaiterControl()

	return &MainContent{
		tableControl:         tableControl,
		waiterControl:        waiterControl,
		stagerManager:        stageManager,
		mainContentContainer: mainContentContainer,
	}
}
