package gui

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

func DefaultNewGameProps() NewGameUI {
	starterPokemonSelectedBtns := make([]*widget.Clickable, len(data.StartPokemonIds))
	for i := range data.StartPokemonIds {
		starterPokemonSelectedBtns[i] = new(widget.Clickable)
	}
	return NewGameUI{
		SelectedIndex:              -1,
		StarterPokemonSelectedBtns: starterPokemonSelectedBtns,
		StartGameBtn:               new(widget.Clickable),
		NewTrainerEditor:           widget.Editor{},
		StarterPokemonList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
}

type NewGameUI struct {
	SelectedIndex              int
	StarterPokemonSelectedBtns []*widget.Clickable
	StartGameBtn               *widget.Clickable

	NewTrainerEditor   widget.Editor
	StarterPokemonList *widget.List
}

func (g *Gui) RenderNewGameScreen(gtx layout.Context) layout.Dimensions {
	// Create a theme for styling

	xInset := gtx.Dp(unit.Dp(200))
	yInset := gtx.Dp(unit.Dp(30))

	return layout.Inset(layout.Inset{
		Top:    unit.Dp(yInset),
		Left:   unit.Dp(xInset),
		Right:  unit.Dp(xInset),
		Bottom: unit.Dp(yInset),
	}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Title
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				title := material.H4(g.Theme, "New Game")
				title.Font.Weight = font.Bold
				return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
			}),

			layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H6(g.Theme, "Enter your Trainer Name")
						title.Font.Weight = font.Bold
						return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// text field to input name
						g.NewGame.NewTrainerEditor.SingleLine = true // Ensure single-line input
						g.NewGame.NewTrainerEditor.MaxLen = 20
						g.NewGame.NewTrainerEditor.Submit = true

						return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Editor(g.Theme, &g.NewGame.NewTrainerEditor, "Enter Name...").Layout(gtx)
						})
					}),
				)
			}),

			layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H6(g.Theme, "Select a Starter Pokemon")
						title.Font.Weight = font.Bold
						return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						starterPokemonsList := g.fetchStarterPokemonList()
						return g.renderStarterPokemonGallery(gtx, starterPokemonsList)
					}),
				)
			}),

			// Start Game button
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return g.renderStartGameBtn(gtx)
			}),
		)
	})
}

func (g *Gui) renderStarterPokemonGallery(gtx layout.Context, starterPokemons []data.BasePokemon) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(g.Theme, g.NewGame.StarterPokemonList).Layout(gtx, len(starterPokemons), func(gtx layout.Context, index int) layout.Dimensions {
				return g.renderStarterPokemonRow(gtx, starterPokemons, index)
			})
		}),
	)
}

func (g *Gui) renderStarterPokemonRow(gtx layout.Context, starterPokemons []data.BasePokemon, index int) layout.Dimensions {
	// Create a row with up to 3 pokemons
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var children []layout.FlexChild
			for i := 0; i < 3; i++ {
				idx := index*3 + i
				if idx >= len(starterPokemons) {
					// Add placeholder for missing trainers
					children = append(children, layout.Flexed(0.33, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}))
					continue
				}
				starterPokemon := starterPokemons[idx]
				children = append(children, layout.Flexed(0.33, func(gtx layout.Context) layout.Dimensions {
					return g.renderStarterPokemonCard(gtx, starterPokemon, idx)
				}))
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderStarterPokemonCard(gtx layout.Context, starterPokemon data.BasePokemon, index int) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(300)))

		isSelected := g.NewGame.SelectedIndex == index
		bgColor := PrimaryBackgroundColor // Default background
		if isSelected {
			bgColor = SelectedColor // Highlighted background
		}

		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
		fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, bgColor, 8)

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			// Trainer Image
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				img := g.loadPokemonImage(starterPokemon.SpritesURL.FrontPath)
				img.Fit = widget.Contain

				// Scale down by applying an inset (adjust Dp as needed)
				inset := layout.UniformInset(unit.Dp(20)) // Reduces size

				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return img.Layout(gtx)
				})
			}),
			// Trainer Name
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(starterPokemon.Name))
				name.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return name.Layout(gtx)
				})
			}),
			// Clickable Overlay for unlocked trainers
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if g.NewGame.StarterPokemonSelectedBtns[index].Clicked(gtx) {
					g.NewGame.SelectedIndex = index
				}
				return g.NewGame.StarterPokemonSelectedBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: cardDims}
				})
			}),
		)
	})

}

func (g *Gui) isStartGameButtonDisabled() bool {
	return g.NewGame.SelectedIndex == -1 || len(g.NewGame.NewTrainerEditor.Text()) < 1
}

func (g *Gui) renderStartGameBtn(gtx layout.Context) layout.Dimensions {

	btnBgColor := SecondaryBackgroundColor
	if g.isStartGameButtonDisabled() {
		btnBgColor = ButtonDisabledColor
	}

	if g.NewGame.StartGameBtn.Clicked(gtx) {
		// start a new game
		g.opts.GameManager = gamestate.NewGameStateManager(
			g.prepareNewGameUser(g.NewGame.NewTrainerEditor.Text(), g.fetchStarterPokemonList()[g.NewGame.SelectedIndex]),
			"", g.NewGame.NewTrainerEditor.Text(),
		)
		g.opts.UserManager = usermanagement.NewUserManager(usermanagement.UserOpts{GameState: g.opts.GameManager})
		// slog.Info("Starting new game", "user", *g.opts.GameManager.Get())
		g.NewGame = DefaultNewGameProps()
		g.SetCurrentScreen(HomeScreen)
	}
	if g.NewGame.StartGameBtn.Hovered() {
		btnBgColor = ButtonHoveredColor
	}

	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		maxDims := gtx.Constraints.Max
		radius := 8

		// define rounded border and fill
		drawRoundedBorder(gtx, maxDims, SecondaryBackgroundColor, unit.Dp(1), unit.Dp(radius))
		fillRoundedShape(gtx, image.Rectangle{Max: maxDims}, btnBgColor, radius)

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				textStyle := material.Body2(g.Theme, "Start Game")
				textStyle.Alignment = text.Middle
				textStyle.Color = TextColor

				return textStyle.Layout(gtx)
			}),
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if g.isStartGameButtonDisabled() {
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}

				return g.NewGame.StartGameBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			}),
		)
	})
}
