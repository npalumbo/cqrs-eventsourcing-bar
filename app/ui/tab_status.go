package ui

import (
	"cqrseventsourcingbar/queries"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TabStatusStage = "TabStatus"

type tabStatusScreen struct {
	tabStatusCard     *widget.Card
	container         *fyne.Container
	tabStatus         *queries.TabStatus
	tableLabel        *widget.Label
	toServeWithAmount *[]tabItemWithAmount
	servedWithAmount  *[]tabItemWithAmount
}

func (t *tabStatusScreen) ExecuteOnTakeOver(param interface{}) {
	tabStatus := param.(*queries.TabStatus)
	t.tabStatus = tabStatus
	t.tabStatusCard.Title = fmt.Sprintf("Tab status for table %d", tabStatus.TableNumber)

	*t.toServeWithAmount = nil
	*t.toServeWithAmount = append(*t.toServeWithAmount, getTabItemsWithAmount(t.tabStatus.ToServe)...)
	*t.servedWithAmount = nil
	*t.servedWithAmount = append(*t.servedWithAmount, getTabItemsWithAmount(t.tabStatus.Served)...)
	t.tableLabel.Text = fmt.Sprintf("%d", tabStatus.TableNumber)
	t.container.Refresh()
}

func (t *tabStatusScreen) GetPaintedContainer() *fyne.Container {
	return t.container
}

func (t *tabStatusScreen) GetStageName() string {
	return TabStatusStage
}

func CreateTabStatusScreen(stageManager *StageManager) *tabStatusScreen {

	tableLabel := widget.NewLabel("")
	toServeWithAmount := &[]tabItemWithAmount{}
	toServeItemList := CreateTabItemList(toServeWithAmount)

	servedWithAmount := &[]tabItemWithAmount{}
	servedItemList := CreateTabItemList(servedWithAmount)

	screenContainer := container.NewStack()

	containerInCard := container.NewBorder(nil, widget.NewButton("Back", func() {
		err := stageManager.TakeOver(MainContentStage, nil)
		if err != nil {
			slog.Error("error launching main content screen", slog.Any("error", err))
		}
	}), nil, nil, container.NewGridWithColumns(2,
		widget.NewCard("", "To Serve", toServeItemList),
		widget.NewCard("", "Served", servedItemList),
	))

	tabStatusCard := widget.NewCard("Tab status", "", containerInCard)
	screenContainer.Add(tabStatusCard)

	return &tabStatusScreen{
		container:         screenContainer,
		tableLabel:        tableLabel,
		toServeWithAmount: toServeWithAmount,
		servedWithAmount:  servedWithAmount,
		tabStatusCard:     tabStatusCard,
	}
}
