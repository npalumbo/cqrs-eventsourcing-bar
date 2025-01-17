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
	waiterDash           *waiterDash
	stagerManager        *StageManager
	mainContentContainer *fyne.Container
}

// ExecuteOnTakeOver implements Stager.
func (m *MainContent) ExecuteOnTakeOver(param interface{}) {
	time.Sleep(200 * time.Millisecond)
	m.tableControl.UpdateActiveTables()
}

// GetPaintedContainer implements Stager.
func (m *MainContent) GetPaintedContainer() *fyne.Container {
	return m.mainContentContainer
}

// GetStageName implements Stager.
func (m *MainContent) GetStageName() string {
	return MainContentStage
}

func CreateMainContent(totalTables int, client *apiclient.ReadClient, stageManager *StageManager, waiters []string, w fyne.Window) *MainContent {
	mainContentContainer := container.NewVBox()

	tableControl := CreateTableControl(totalTables, client, waiters, stageManager, mainContentContainer)

	waiterDash := CreateWaiterDash(mainContentContainer)

	tableControl.UpdateActiveTables()

	return &MainContent{
		tableControl:         tableControl,
		waiterDash:           waiterDash,
		stagerManager:        stageManager,
		mainContentContainer: mainContentContainer,
	}
}
