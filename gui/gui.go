package gui

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
)

type GuiOpts struct {
	UserService    usermanagement.UserManager
	PokemonService pokemon.PokemonManager

	ActionContainer *fyne.Container
	LogContainer    *fyne.Container
	BattleLogChan   <-chan string
}

func InitializeGUI(opts GuiOpts) {
	a := app.New()

	userWindow := a.NewWindow("PokéSim")

	user, err := opts.UserService.LoadUser()
	if err != nil {
		slog.Error("Failed to load user", "error", err)
	}
	// slog.Info("User loaded", "user", user)
	opts.ActionContainer = container.NewStack()
	opts.LogContainer = container.NewStack()

	// Create the main action area
	opts.ActionContainer = opts.DefaultActionTab()

	// Create the log area
	opts.LogContainer = opts.NewLogContainer()
	logContainer := container.NewVScroll((opts.LogContainer))
	logContainer.SetMinSize(fyne.NewSize(200, 300))

	// Create the trainer sidebar
	trainerSideBar := opts.GetTrainerSideBar(user)

	// Set up the layout
	content := container.New(layout.NewBorderLayout(nil, nil, nil, trainerSideBar),
		trainerSideBar, // Right
		container.NewBorder(nil, 
			addBorder(logContainer), 
			nil, nil, opts.ActionContainer),
	)

	userWindow.SetContent(content)
	userWindow.SetMainMenu(opts.MainMenu())

	userWindow.Resize(fyne.NewSize(1200, 900))
	userWindow.SetFixedSize(true) // Disallow resizing by the user
	userWindow.ShowAndRun()
}

// addBorder adds padding and a border to a given content area
func addBorder(content fyne.CanvasObject) fyne.CanvasObject {
	// Create a padded container for margin-like space
	paddedContent := container.NewPadded(content)

	// Create a rectangle to serve as the border
	border := canvas.NewRectangle(color.Black)
	border.StrokeWidth = 2
	border.StrokeColor = color.RGBA{0, 0, 0, 255}
	border.FillColor = color.Transparent

	// Overlay the border and content using a Border layout
	borderContainer := container.New(layout.NewBorderLayout(nil, nil, nil, nil), border, paddedContent)

	return borderContainer
}
