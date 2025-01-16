package ui

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tableButton struct {
	widget.Button
	ID           int
	menuActive   *fyne.Menu
	menuInactive *fyne.Menu
	Active       bool
	waiters      []string
	stageManager *StageManager
}

func newTableButton(ID int, waiters []string, stageManager *StageManager) *tableButton {
	tableButton := &tableButton{ID: ID, waiters: waiters, stageManager: stageManager}
	tableButton.ExtendBaseWidget(tableButton)

	tableButton.menuActive = fyne.NewMenu("Active Table",
		fyne.NewMenuItem("Close Table", func() { fmt.Println("Clicked Close") }),
	)
	tableButton.menuInactive = fyne.NewMenu("Active Table",
		fyne.NewMenuItem("Open Tab", func() {
			err := stageManager.TakeOver(OpenTabStage, ID)
			slog.Error("error launching open tab screen", slog.Any("error", err))
		}),
	)
	tableButton.Text = fmt.Sprintf("Table %d", ID)
	return tableButton

}

func (t *tableButton) Tapped(e *fyne.PointEvent) {
	var menu *fyne.Menu
	if t.Active {
		menu = t.menuActive
	} else {
		menu = t.menuInactive
	}
	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(t), e.AbsolutePosition)
}

func (t *tableButton) TappedSecondary(_ *fyne.PointEvent) {
}

func (t *tableButton) SetActive() {
	t.Active = true
	t.Importance = widget.HighImportance
}

func (t *tableButton) SetInactive() {
	t.Active = false
	t.Importance = widget.LowImportance
}
