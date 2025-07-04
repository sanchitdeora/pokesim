package gui

import (
	"image"
	"log/slog"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/utils"
)

func (g *Gui) RenderWildScreen(gtx layout.Context) layout.Dimensions {
	// Use a vertical layout for screen elements
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Wild Environment")
			title.Font.Weight = font.Bold
			return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					environments := g.WildEnvironmentProps.Environments
					return g.renderEnvironmentsGallery(gtx, environments)
				}))
		}),
	)
}

func (g *Gui) renderEnvironmentsGallery(gtx layout.Context, environments []EnvironmentUI) layout.Dimensions {
	// Wrap the environment list in a flex layout for better control

	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(g.Theme, g.WildEnvironmentProps.WildEnvironmentList).Layout(gtx, len(environments), func(gtx layout.Context, index int) layout.Dimensions {
				return g.renderEnvironmentRow(gtx, environments, index)
			})
		}),
	)
}

func (g *Gui) renderEnvironmentRow(gtx layout.Context, environments []EnvironmentUI, index int) layout.Dimensions {
	// Create a row with up to 4 environments
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var children []layout.FlexChild
			for i := 0; i < 4; i++ {
				idx := index*4 + i
				if idx >= len(environments) {
					// Add placeholder for missing environments
					children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}))
					continue
				}
				environment := environments[idx]
				children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
					return g.renderEnvironmentCard(gtx, environment)
				}))
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderEnvironmentCard(gtx layout.Context, environment EnvironmentUI) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(175)))

		// Paint the rounded border
		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			// Environment Image
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				newCtx := gtx
				newCtx.Constraints.Max = cardDims
				newCtx.Constraints.Max.Y -= 50
				
				img := environment.Image
				img.Fit = widget.Fill
				img.Position = layout.Center
				return img.Layout(newCtx)
			}),
			// Environment Name
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(string(environment.Name)))
				name.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return name.Layout(gtx)
				})
			}),
			// Clickable Overlay for Environments
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if environment.EnvironmentBtn.Clicked(gtx) {
					slog.Info("clicked on environment", "name", environment.Name)
					g.Battle = g.NewWildBattle(g.opts.UserManager.GetUser(), g.opts.PokemonService.SearchWildPokemon(environment.Name))
					return g.LoadBattle(gtx)
				}

				return environment.EnvironmentBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: cardDims}
				})
			}),
		)
	})
}

