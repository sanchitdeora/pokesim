package gui

import (
	"image"
	"image/color"
	"log"
	"log/slog"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

var (
	PrimaryBackgroundColor   = color.NRGBA{R: 248, G: 248, B: 251, A: 255}
	SecondaryBackgroundColor = color.NRGBA{R: 232, G: 232, B: 243, A: 255}
	SelectedColor            = color.NRGBA{R: 203, G: 203, B: 212, A: 255}

	ButtonDisabledColor = color.NRGBA{R: 232, G: 232, B: 243, A: 63}
	ButtonHoveredColor  = color.NRGBA{R: 203, G: 203, B: 212, A: 255}

	TextColor          = color.NRGBA{R: 14, G: 14, B: 27, A: 255}
	SecondaryTextColor = color.NRGBA{R: 80, G: 79, B: 150, A: 255}

	CardBorderColor = color.NRGBA{R: 209, G: 208, B: 230, A: 255}

	RedBtnColor = color.NRGBA{R: 200, G: 0, B: 0, A: 255}
)

type Screen int

const (
	LoadGameScreen Screen = iota
	NewGameScreen
	HomeScreen
	BattleScreen
	TrainerScreen
	WildScreen
	PartyScreen
	PokemonSummaryScreen
	BoxScreen
	BagScreen
	ShopScreen
	ToBeReplaced
)

type GuiOpts struct {
	GameManager    gamestate.GameStateManager
	UserManager    usermanagement.UserManager
	PokemonService pokemon.PokemonService
	BoxManager     pokemon.BoxManager
}

type Gui struct {
	opts          GuiOpts
	CurrentWindow Screen
	Theme         *material.Theme
	Buttons       map[Screen]*widget.Clickable

	// Trainer/Wild and Battle Props
	TrainerList          *widget.List
	WildEnvironmentProps WildEnvironmentProps
	Battle               Battle

	// Pokemon Party and Summary Props
	Party          PartyProps
	SummaryPokemon *data.Pokemon

	NewGame   NewGameUI
	LoadGames []LoadGameUI

	BoxList *widget.List
	Box     BoxProps

	// Shop Props
	Shop     ShopProps
	ItemList *widget.List
}

func InitializeGUI() {
	go func() {
		// Set up the window
		window := setWindowOptions()

		// Set up theme
		theme := setTheme()

		// State management
		buttons := map[Screen]*widget.Clickable{
			LoadGameScreen:       new(widget.Clickable),
			NewGameScreen:        new(widget.Clickable),
			HomeScreen:           new(widget.Clickable),
			TrainerScreen:        new(widget.Clickable),
			WildScreen:           new(widget.Clickable),
			BattleScreen:         new(widget.Clickable),
			PartyScreen:          new(widget.Clickable),
			PokemonSummaryScreen: new(widget.Clickable),
			BoxScreen:            new(widget.Clickable),
			BagScreen:            new(widget.Clickable),
			ShopScreen:           new(widget.Clickable),
			ToBeReplaced:         new(widget.Clickable),
		}

		opts := GuiOpts{
			PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
		}

		gui := &Gui{
			opts:          opts,
			CurrentWindow: LoadGameScreen,
			Theme:         theme,
			Buttons:       buttons,

			TrainerList: &widget.List{
				List: layout.List{Axis: layout.Vertical},
			},

			// WildEnvironmentProps: DefaultWildEnvironmentProps(),

			// Party:          DefaultPartyProps(),
			SummaryPokemon: nil,

			BoxList: &widget.List{
				List: layout.List{Axis: layout.Vertical},
			},

			NewGame:   NewGameUI{},
			LoadGames: createLoadGamesUI(),

			// Shop: DefaultShopProps(),
			ItemList: &widget.List{
				List: layout.List{Axis: layout.Vertical},
			},
		}

		if err := run(window, gui); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, gui *Gui) error {
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			paint.Fill(gtx.Ops, SecondaryBackgroundColor)
			gui.drawCentralContainer(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func (g *Gui) drawCentralContainer(gtx layout.Context) layout.Dimensions {
	// Container size in Dp
	containerWidth := gtx.Dp(unit.Dp(1400))
	containerHeight := gtx.Dp(unit.Dp(800))

	// slog.Info("Container Size", "Width", containerWidth, "Height", containerHeight, "Max Width", gtx.Constraints.Max.X, "Max Height", gtx.Constraints.Max.Y)

	// Center the container
	xOffset := (gtx.Constraints.Max.X - containerWidth) / 2
	yOffset := (gtx.Constraints.Max.Y - containerHeight) / 2

	// Apply the offset to the layout (move the content to the center)
	ts := op.Offset(image.Pt(xOffset, yOffset)).Push(gtx.Ops)
	defer ts.Pop() // Remove the offset after we're done with the layout

	// Clip and draw the container with rounded corners
	defer clip.RRect{
		Rect: image.Rectangle{
			Min: image.Pt(0, 0), // Position it at the offset (0,0 after offset is applied)
			Max: image.Pt(containerWidth, containerHeight),
		},
		SE: 15, NW: 15, NE: 15, SW: 15, // Corner radii
	}.Push(gtx.Ops).Pop()

	// Fill the central container background (white)
	paint.Fill(gtx.Ops, PrimaryBackgroundColor)

	// Define insets to apply padding inside the container
	insets := layout.Inset{
		Top:    unit.Dp(10),
		Bottom: unit.Dp(50),
		Left:   unit.Dp(50),
		Right:  unit.Dp(50),
	}

	gtx.Constraints.Min = image.Pt(gtx.Dp(unit.Dp(0)), gtx.Dp(unit.Dp(0)))
	gtx.Constraints.Max = image.Pt(containerWidth, containerHeight)

	// Apply padding and layout the content inside the central container
	return insets.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceBetween, // Add space between menu and content
		}.Layout(gtx,
			// Menu Bar at the top (place it at the top of the flex layout)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.RenderMenu(gtx)
			}),

			// Separator line
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Temporarily offset the drawing context to the leftmost edge
				offset := op.Offset(image.Pt(-gtx.Dp(unit.Dp(50)), 0)).Push(gtx.Ops) // Remove the left inset
				defer offset.Pop()                                                   // Restore the original offset after drawing

				// Define the separator dimensions
				separatorWidth := gtx.Constraints.Max.X + gtx.Dp(unit.Dp(50))*2 // Account for both insets
				separator := layout.Dimensions{Size: image.Point{X: separatorWidth, Y: 2}}

				// Draw the separator line
				fillRoundedShape(gtx, image.Rectangle{Max: separator.Size}, SecondaryBackgroundColor, 0)

				return separator
			}),

			// Main content area below the menu bar
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				// Content of the main white container
				xInset := gtx.Dp(unit.Dp(200))
				yInset := gtx.Dp(unit.Dp(30))

				return layout.Inset(layout.Inset{
					Top:    unit.Dp(yInset),
					Left:   unit.Dp(xInset),
					Right:  unit.Dp(xInset),
					Bottom: unit.Dp(yInset),
				}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.RenderCurrentScreen(gtx)
				})
			}),
		)
	})
}

func loadCustomFont() *text.Shaper {
	// initialze custom font
	ManropeFontPath := utils.GetFullPath("assets/fonts/Manrope-SemiBold.ttf")
	ManropeBoldFontPath := utils.GetFullPath("assets/fonts/Manrope-Bold.ttf")
	ManropeExtraBoldFontPath := utils.GetFullPath("assets/fonts/Manrope-ExtraBold.ttf")

	// Create a new Gio text font collection with the custom font
	fontCollection := []text.FontFace{
		{
			Font: font.Font{
				Weight: font.Normal, // You can adjust weight (e.g., Bold, Medium)
			},
			Face: getCustomFontFace(ManropeFontPath),
		},
		{
			Font: font.Font{
				Weight: font.Bold, // You can adjust weight (e.g., Bold, Medium)
			},
			Face: getCustomFontFace(ManropeBoldFontPath),
		},
		{
			Font: font.Font{
				Weight: font.ExtraBold, // You can adjust weight (e.g., Bold, Medium)
			},
			Face: getCustomFontFace(ManropeExtraBoldFontPath),
		},
	}

	// Return a new Material theme with the custom font
	return text.NewShaper(text.WithCollection(fontCollection))
}

func setWindowOptions() *app.Window {
	w := new(app.Window)
	w.Option(app.Title("PokéSim"))
	w.Option(app.Maximized.Option())
	// window.Option(app.Decorated(false))

	return w
}

func setTheme() *material.Theme {
	t := material.NewTheme()
	t.Shaper = loadCustomFont()
	t.Bg = PrimaryBackgroundColor

	return t
}

func getCustomFontFace(filepath string) opentype.Face {
	// Read the font file
	fontBytes, err := os.ReadFile(filepath)
	if err != nil {
		slog.Error("Failed to read font file: %v", "error", err)
	}

	// Parse the fb
	fb, err := opentype.Parse(fontBytes)
	if err != nil {
		slog.Error("Failed to parse font: %v", "error", err)
	}

	return fb
}
