package main

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/theme"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Game with Navigation Bar")

// 	// Create toolbar items for each action
// 	newGameButton := widget.NewToolbarAction(theme.DeleteIcon(), func() {
// 		// Add your "New Game" logic here
// 	})
// 	loadGameButton := widget.NewToolbarAction(theme.DeleteIcon(), func() {
// 		// Add your "Load Game" logic here
// 	})
// 	saveGameButton := widget.NewToolbarAction(theme.DeleteIcon(), func() {
// 		// Add your "Save Game" logic here
// 	})
// 	settingsButton := widget.NewToolbarAction(theme.DeleteIcon(), func() {
// 		// Add your "Settings" logic here
// 	})

// 	// Create the toolbar using the buttons
// 	toolbar := widget.NewToolbar(
// 		widget.NewToolbarAction(theme.DeleteIcon(), func() {
// 			// New Game Logic
// 		}),
// 		newGameButton,
// 		loadGameButton,
// 		saveGameButton,
// 		settingsButton,
// 	)

// 	// Main content area, for example, an action display
// 	actionContent := widget.NewLabel("Welcome to the Game!")

// 	// Layout with toolbar on top and main content below
// 	content := container.NewBorder(toolbar, nil, nil, nil, actionContent)

// 	// Set the content and run the app
// 	myWindow.SetContent(content)
// 	myWindow.Resize(fyne.NewSize(800, 600))
// 	myWindow.ShowAndRun()
// }
