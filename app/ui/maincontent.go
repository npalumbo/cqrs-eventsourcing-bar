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
	stagerManager        *StageManager
	mainContentContainer *fyne.Container
}

func (m *MainContent) ExecuteOnTakeOver(param interface{}) {
	time.Sleep(200 * time.Millisecond)
	m.tableControl.UpdateActiveTables()
}

func (m *MainContent) GetPaintedContainer() *fyne.Container {
	return m.mainContentContainer
}

func (m *MainContent) GetStageName() string {
	return MainContentStage
}

func CreateMainContent(totalTables int, client *apiclient.ReadClient, stageManager *StageManager, waiters []string, w fyne.Window) *MainContent {
	mainContentContainer := container.NewStack()

	tableControl := CreateTableControl(totalTables, client, waiters, stageManager, mainContentContainer)

	tableControl.UpdateActiveTables()

	return &MainContent{
		tableControl:         tableControl,
		stagerManager:        stageManager,
		mainContentContainer: mainContentContainer,
	}
}
