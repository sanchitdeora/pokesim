package gui

import (
	"log/slog"

	"gioui.org/font"
	"gioui.org/layout"
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
}
