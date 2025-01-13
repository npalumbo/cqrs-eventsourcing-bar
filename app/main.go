package main

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/readservice/model"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("CQRS ES BAR")

	createTableButton := func(i int) *widget.Button {

		button := widget.NewButton(fmt.Sprintf("Mesa %d", i), func() {
			println("Circle button clicked!")
		})
		return button
	}

	url := "http://localhost:8081/activeTableNumbers" // Replace with the actual URL

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP GET request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check for response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP status code error: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// Parse the JSON response
	var active model.QueryResponse[model.ActiveTableNumbersResponse]
	err = json.Unmarshal(body, &active)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Create a slice to hold the six buttons
	buttons := make([]*widget.Button, 6)

	// Create the six buttons
	for i := 0; i < 6; i++ {
		buttons[i] = createTableButton(i)
	}

	for i := 0; i < len(active.Data.ActiveTables); i++ {
		button := buttons[i]
		button.Importance = 1
	}

	// Arrange the buttons using a grid layout
	grid := container.New(layout.NewGridLayout(3), buttons[0], buttons[1], buttons[2], buttons[3], buttons[4], buttons[5])

	w.SetContent(grid)
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
