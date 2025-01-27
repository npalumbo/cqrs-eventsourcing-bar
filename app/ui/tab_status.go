package ui

import (
	"golangsevillabar/queries"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TabStatusStage = "TabStatus"

type tabStatusScreen struct {
	container *fyne.Container
	tabStatus *queries.TabStatus
	form      *widget.Form
}

func (t *tabStatusScreen) ExecuteOnTakeOver(param interface{}) {
	tabStatus := param.(*queries.TabStatus)
	t.tabStatus = tabStatus
}

func (t *tabStatusScreen) GetPaintedContainer() *fyne.Container {
	return t.container
}

func (t *tabStatusScreen) GetStageName() string {
	return TabStatusStage
}

func CreateTabStatusCreen() *tabStatusScreen {
	container := container.NewVBox()
	form := widget.NewForm()

	container.Add(widget.NewCard("Tab status", "", form))

	tabStatusScreen := &tabStatusScreen{container: container, form: form}

	return tabStatusScreen
}
