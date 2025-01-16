package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type waiterDash struct {
	container *fyne.Container
}

func CreateWaiterDash(mainContentContainer *fyne.Container) *waiterDash {
	waiterDashContainer := container.NewVBox(widget.NewCard("Waiters", "", widget.NewForm()))
	mainContentContainer.Add(waiterDashContainer)
	return &waiterDash{
		container: waiterDashContainer,
	}
}
