package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type StageManager struct {
	currentViewContainer *fyne.Container
	stagerMap            map[string]Stager
}

type StagerController interface {
	TakeOver(name string, param interface{}) error
	RegisterStager(stager Stager)
	GetContainer() *fyne.Container
}

type DefaultStager struct {
}

type Stager interface {
	GetPaintedContainer() *fyne.Container
	ExecuteOnTakeOver(param interface{})
	GetStageName() string
}

func CreateStageManager() StageManager {
	return StageManager{
		currentViewContainer: container.NewStack(),
		stagerMap:            make(map[string]Stager),
	}
}

func (s StageManager) RegisterStager(stager Stager) {
	s.stagerMap[stager.GetStageName()] = stager
}

func (s StageManager) TakeOver(name string, param interface{}) error {
	stager, ok := s.stagerMap[name]

	if !ok {
		return errors.New("Unknown stager: " + name)
	}

	s.currentViewContainer.RemoveAll()
	container := stager.GetPaintedContainer()
	container.Refresh()
	s.currentViewContainer.Add(container)
	stager.ExecuteOnTakeOver(param)
	s.currentViewContainer.Refresh()

	return nil
}

func (s StageManager) GetContainer() *fyne.Container {
	return s.currentViewContainer
}

func (d *DefaultStager) ExecuteOnTakeOver() {

}
