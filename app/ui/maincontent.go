package ui

import (
	"golangsevillabar/app/apiclient"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent struct {
	tableControl         *tableControl
	stagerManager        StageManager
	mainContentContainer *fyne.Container
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return m.mainContentContainer
}

func CreateMainContent(totalTables int, client *apiclient.Client, w fyne.Window) *MainContent {
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)

	mainContentContainer := container.NewVBox()

	tableControl := CreateTableControl(totalTables, client, mainContentContainer)
	mainContentContainer.Add(currentContainer)

	tableControl.UpdateActiveTables()

	return &MainContent{
		tableControl:         tableControl,
		stagerManager:        stageManager,
		mainContentContainer: mainContentContainer,
	}
}
