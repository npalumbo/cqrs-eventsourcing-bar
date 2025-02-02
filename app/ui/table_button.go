package ui

import (
	"cqrseventsourcingbar/queries"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tableNumberAndTabId struct {
	tableNumber int
	tabId       string
}

type tableButton struct {
	widget.Button
	ID           int
	menuActive   *fyne.Menu
	menuInactive *fyne.Menu
	Active       bool
	stageManager *StageManager
	tabStatus    *queries.TabStatus
}

func newTableButton(ID int, stageManager *StageManager) *tableButton {
	tableButton := &tableButton{ID: ID, stageManager: stageManager}
	tableButton.ExtendBaseWidget(tableButton)

	tableButton.menuActive = fyne.NewMenu("Active Table",
		fyne.NewMenuItem("Invoice for Table", func() {
			err := stageManager.TakeOver(InvoiceStage, ID)
			if err != nil {
				slog.Error("error launching invoice screen", slog.Any("error", err))
			}
		}),
		fyne.NewMenuItem("Order drinks", func() {
			err := stageManager.TakeOver(PlaceOrderStage, tableNumberAndTabId{
				tableNumber: ID,
				tabId:       tableButton.tabStatus.TabID,
			})
			if err != nil {
				slog.Error("error launching order drinks screen", slog.Any("error", err))
			}
		}),
		fyne.NewMenuItem("Tab status", func() {
			err := stageManager.TakeOver(TabStatusStage, tableButton.tabStatus)
			if err != nil {
				slog.Error("error launching tab status screen", slog.Any("error", err))
			}
		}),
	)
	tableButton.menuInactive = fyne.NewMenu("Inactive Table",
		fyne.NewMenuItem("Open Tab", func() {
			err := stageManager.TakeOver(OpenTabStage, ID)
			if err != nil {
				slog.Error("error launching open tab screen", slog.Any("error", err))
			}
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

func (t *tableButton) SetActive(tabStatus *queries.TabStatus) {
	t.Active = true
	t.Importance = widget.HighImportance
	t.tabStatus = tabStatus
}

func (t *tableButton) SetInactive() {
	t.Active = false
	t.Importance = widget.LowImportance
}
