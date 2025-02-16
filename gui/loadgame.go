package gui

import (
	"image"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/usermanagement"
)

type LoadGameUI struct {
	LoadGameButton *widget.Clickable
	FileName       string
	Image          widget.Image
}

func (g *Gui) RenderLoadGameScreen(gtx layout.Context) layout.Dimensions {
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
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(g.Theme, "Load Game")
				title.Font.Weight = font.Bold
				return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						loadGameList := g.createLoadGamesUI()
						return g.renderLoadGameGallery(gtx, loadGameList)
					}))
			}),
		)
	})
}

func (g *Gui) renderLoadGameGallery(gtx layout.Context, loadGameList []LoadGameUI) layout.Dimensions {
	// Wrap the trainer list in a flex layout for better control

	gameList := &widget.List{
		List: layout.List{Axis: layout.Vertical},
	}

	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(g.Theme, gameList).Layout(gtx, len(loadGameList), func(gtx layout.Context, index int) layout.Dimensions {
				return g.renderLoadGameRow(gtx, loadGameList, index)
			})
		}),
	)
}

func (g *Gui) renderLoadGameRow(gtx layout.Context, loadGameList []LoadGameUI, index int) layout.Dimensions {
	// Create a row with up to 4 trainers
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var children []layout.FlexChild
			for i := 0; i < 4; i++ {
				idx := index*4 + i
				if idx >= len(loadGameList) {
					// Add placeholder for missing trainers
					children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}))
					continue
				}
				loadGame := loadGameList[idx]
				children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
					return g.renderLoadGameCard(gtx, loadGame)
				}))
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderLoadGameCard(gtx layout.Context, loadGame LoadGameUI) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(250)))

		// Paint the rounded border
		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			// Trainer Image
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				img := loadGame.Image
				img.Fit = widget.Contain
			
				// Scale down by applying an inset (adjust Dp as needed)
				inset := layout.UniformInset(unit.Dp(20)) // Reduces size
			
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {			
					return img.Layout(gtx)
				})
			}),
			// Trainer Name
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				name := material.Body1(g.Theme, displayFileName(loadGame.FileName))
				name.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return name.Layout(gtx)
				})
			}),
			// Clickable Overlay for unlocked trainers
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if loadGame.FileName != "New Game" {
					if g.Buttons[LoadGameScreen].Clicked(gtx) {
						g.opts.GameManager = gamestate.NewGameStateManager(nil, "", displayFileName(loadGame.FileName))
						g.opts.UserManager = usermanagement.NewUserManager(usermanagement.UserOpts{GameState: g.opts.GameManager})
						g.SetCurrentScreen(HomeScreen)
					}
					return g.Buttons[LoadGameScreen].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: cardDims}
					})
				} else {
					if g.Buttons[NewGameScreen].Clicked(gtx) {
						g.NewGame = DefaultNewGameProps()
						g.SetCurrentScreen(NewGameScreen)
					}
					return g.Buttons[NewGameScreen].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: cardDims}
					})
				}
			}),
		)
	})
}

func displayFileName(fileName string) string {
	return strings.TrimSuffix(fileName, ".json")
}