package ui

import (
	"fmt"
	"golangsevillabar/queries"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TabStatusStage = "TabStatus"

type tabStatusScreen struct {
	container         *fyne.Container
	tabStatus         *queries.TabStatus
	tableLabel        *widget.Label
	toServeWithAmount *[]tabItemWithAmount
	servedWithAmount  *[]tabItemWithAmount
}

func (t *tabStatusScreen) ExecuteOnTakeOver(param interface{}) {
	tabStatus := param.(*queries.TabStatus)
	t.tabStatus = tabStatus

	*t.toServeWithAmount = nil
	*t.toServeWithAmount = append(*t.toServeWithAmount, getTabItemsWithAmount(t.tabStatus.ToServe)...)
	*t.servedWithAmount = nil
	*t.servedWithAmount = append(*t.servedWithAmount, getTabItemsWithAmount(t.tabStatus.Served)...)
	t.tableLabel.Text = fmt.Sprintf("%d", tabStatus.TableNumber)
}

func (t *tabStatusScreen) GetPaintedContainer() *fyne.Container {
	return t.container
}

func (t *tabStatusScreen) GetStageName() string {
	return TabStatusStage
}

func CreateTabStatusCreen(stageManager *StageManager) *tabStatusScreen {

	tableLabel := widget.NewLabel("")
	toServeWithAmount := &[]tabItemWithAmount{}
	toServeItemList := CreateTabItemList(toServeWithAmount)

	servedWithAmount := &[]tabItemWithAmount{}
	servedItemList := CreateTabItemList(servedWithAmount)

	screenContainer := container.NewStack()

	containerInCard := container.NewGridWithColumns(2,
		widget.NewLabel("Table"), tableLabel,
		widget.NewLabel("To Serve"), toServeItemList,
		widget.NewLabel("Served"), servedItemList,
		widget.NewButton("Back", func() {
			err := stageManager.TakeOver(MainContentStage, nil)
			if err != nil {
				slog.Error("error launching main content screen", slog.Any("error", err))
			}
		}),
	)

	screenContainer.Add(widget.NewCard("Tab status", "", containerInCard))

	tabStatusScreen := &tabStatusScreen{container: screenContainer, tableLabel: tableLabel, toServeWithAmount: toServeWithAmount, servedWithAmount: servedWithAmount}

	return tabStatusScreen
}
