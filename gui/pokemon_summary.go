package gui

import (
	"log/slog"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func (g *Gui) SetSummaryPokemonActive(p *data.Pokemon) {
	g.SummaryPokemon = p
	g.SetCurrentScreen(PokemonSummaryScreen)
}

func (g *Gui) RenderPokemonSummary(gtx layout.Context) layout.Dimensions {
	if g.SummaryPokemon == nil {
		slog.Error("No pokemon to display summary")
		g.SetCurrentScreen(HomeScreen)
	}

	// Create a theme for styling
	xInset := gtx.Dp(unit.Dp(200))
	yInset := gtx.Dp(unit.Dp(30))

	return layout.Inset(layout.Inset{
		Top:    unit.Dp(yInset),
		Left:   unit.Dp(xInset),
		Right:  unit.Dp(xInset),
		Bottom: unit.Dp(yInset),
	}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,

			// Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(g.Theme, "Pokemon Summary: "+utils.ToCapitalizeFirstLetterOfEachWord(g.SummaryPokemon.Name))
				title.Font.Weight = font.Bold

				return title.Layout(gtx)
			}),
		)
	})
}
