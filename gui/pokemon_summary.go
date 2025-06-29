package gui

import (
	"fmt"
	"log/slog"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
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
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, fmt.Sprintf("Pokemon Summary: #%04d - %s", g.SummaryPokemon.ID, utils.ToCapitalizeFirstLetterOfEachWord(g.SummaryPokemon.Name)))
			title.Font.Weight = font.Bold

			return title.Layout(gtx)
		}),
		layout.Flexed(0.9, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderPokemonSummary(gtx, *g.SummaryPokemon)
					}),
				)
			})
		}),
	)
}

func (g *Gui) renderPokemonSummary(gtx layout.Context, p data.Pokemon) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Photo
				layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))

						img := loadImage(p.SpritesURL.FrontPath) // Replace with your image loading logic
						img.Fit = widget.Contain
						img.Position = layout.Center
						return (layout.Center.Layout(gtx, img.Layout))
					})
				}),

				// Details
				layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))

						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Pokédex No.: %v", g.SummaryPokemon.ID))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Name: %s", utils.ToCapitalizeFirstLetterOfEachWord(g.SummaryPokemon.Name)))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Level: %v", g.SummaryPokemon.Level))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Type: %s", getPokemonSummaryTypes(*g.SummaryPokemon)))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Exp Points To Next Lvl: %v", g.SummaryPokemon.ExperienceLeft))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return g.renderEXPBar(gtx, g.GetRemainingEXPFraction(g.SummaryPokemon))
									})
								}),
							)
						})
					})
				}),
			)
		}),

		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Stats
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))

						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {

							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("HP: %v", g.SummaryPokemon.Stats.HP.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Attack: %v", g.SummaryPokemon.Stats.Attack.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Defence: %v", g.SummaryPokemon.Stats.Defense.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Sp. Attack: %v", g.SummaryPokemon.Stats.SpecialAttack.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Sp. Defence: %v", g.SummaryPokemon.Stats.SpecialDefense.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
								layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.Theme, fmt.Sprintf("Speed: %v", g.SummaryPokemon.Stats.Speed.Value))
									label.Alignment = text.Middle
									label.Font.Weight = font.Bold
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, label.Layout)
								}),
							)
						})
					})
				}),

				// Moves
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))

						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
									if p.Moveset.Move1 == nil {
										return layout.Dimensions{}
									}
									return g.renderPokemonSummaryMoveTile(gtx, *p.Moveset.Move1)

								}),
								layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
									if p.Moveset.Move2 == nil {
										return layout.Dimensions{}
									}
									return g.renderPokemonSummaryMoveTile(gtx, *p.Moveset.Move2)

								}),
								layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
									if p.Moveset.Move3 == nil {
										return layout.Dimensions{}
									}
									return g.renderPokemonSummaryMoveTile(gtx, *p.Moveset.Move3)

								}),
								layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
									if p.Moveset.Move4 == nil {
										return layout.Dimensions{}
									}
									return g.renderPokemonSummaryMoveTile(gtx, *p.Moveset.Move4)
								}),
							)
						})
					})
				}),
			)
		}),
	)
}

func (g *Gui) renderPokemonSummaryMoveTile(gtx layout.Context, move data.Moves) layout.Dimensions {

	return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))

		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Horizontal,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(move.Name))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return layout.UniformInset(unit.Dp(1)).Layout(gtx, label.Layout)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(string(move.Type)))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return layout.UniformInset(unit.Dp(1)).Layout(gtx, label.Layout)
					}),
				)
			}),

			layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Power: %v", move.Power))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return layout.UniformInset(unit.Dp(1)).Layout(gtx, label.Layout)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Accuracy: %v", move.Accuracy))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return layout.UniformInset(unit.Dp(1)).Layout(gtx, label.Layout)
					}),
				)
			}),
		)
	})
}

func getPokemonSummaryTypes(pokemon data.Pokemon) string {
	if pokemon.Type2 == "" {
		return utils.ToCapitalizeFirstLetterOfEachWord(string(pokemon.Type1))
	}
	return utils.ToCapitalizeFirstLetterOfEachWord(fmt.Sprintf("%s | %s", pokemon.Type1, pokemon.Type2))
}
