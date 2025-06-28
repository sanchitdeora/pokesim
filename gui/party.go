package gui

import (
	"fmt"
	"image"
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

func DefaultPartyProps() PartyProps {
	pokmemonSelectionBtns := make([]*widget.Clickable, 6)
	for i := range pokmemonSelectionBtns {
		pokmemonSelectionBtns[i] = new(widget.Clickable)
	}
	return PartyProps{
		SelectedIndex:       -1,
		PokemonSelectedBtns: pokmemonSelectionBtns,
		SummaryBtn:          new(widget.Clickable),
		MoveUpBtn:           new(widget.Clickable),
		MoveDownBtn:         new(widget.Clickable),
	}
}

type PartyProps struct {
	SelectedIndex       int
	PokemonSelectedBtns []*widget.Clickable
	SummaryBtn          *widget.Clickable
	MoveUpBtn           *widget.Clickable
	MoveDownBtn         *widget.Clickable
}

func (g *Gui) RenderPartyScreen(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,

		// Title
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Party")
			title.Font.Weight = font.Bold

			return title.Layout(gtx)
		}),

		layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset(layout.Inset{
				Top:    unit.Dp(8),
				Left:   unit.Dp(32),
				Right:  unit.Dp(128),
				Bottom: unit.Dp(8),
			}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderPartyList(gtx)
					}))
			})
		}),

		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return g.renderPartyButtonRow(gtx)
		}),
	)
}

func (g *Gui) renderPartyList(gtx layout.Context) layout.Dimensions {
	// Render the Pokémon party page with 6 fixed rows
	party := g.opts.UserManager.GetUser().Party

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var rows []layout.FlexChild
			for i := range 6 {
				index := i
				if i < len(party) {
					// Render the Pokémon in the party
					// slog.Info("in loop", "index", index, "selectedIndex", g.Party.SelectedIndex, "gtx", gtx.Constraints)
					pokemon := party[i]
					rows = append(rows, layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
						if g.Party.PokemonSelectedBtns[index].Clicked(gtx) {
							g.Party.SelectedIndex = index
							slog.Info("pokemon Selected Clicked", "index", index, "selectedIndex", g.Party.SelectedIndex, "gtx", gtx.Constraints)
						}
						return g.renderPartyRow(gtx, pokemon, index)
					}))
				} else {
					// Render an empty row as a placeholder
					rows = append(rows, layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
						return g.renderEmptyPartyRow(gtx)
					}))
				}
			}
			return rows
		}(gtx)...,
	)
}

func (g *Gui) renderPartyRow(gtx layout.Context, pokemon *data.Pokemon, index int) layout.Dimensions {
	if pokemon == nil {
		// Render an empty placeholder without a border
		return layout.Dimensions{}
	}

	isSelected := g.Party.SelectedIndex == index
	bgColor := SecondaryBackgroundColor // Default background
	if isSelected {
		bgColor = SelectedColor // Highlighted background
	}

	// Render a single row for a Pokémon in the party
	return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := gtx.Constraints.Max

		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
		fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, bgColor, 8)

		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return g.renderPartyRowDetail(gtx, pokemon)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return g.Party.PokemonSelectedBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: cardDims}
				})

			}),
		)
	})
}

func (g *Gui) renderPartyButtonRow(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),

		// Summary Btn
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			isBtnDisabled := g.Party.SelectedIndex == -1

			btnBgColor := SecondaryBackgroundColor // Grey background for disabled
			if isBtnDisabled {
				btnBgColor = ButtonDisabledColor // Active background (blue)
			}

			if g.Party.SummaryBtn.Clicked(gtx) {
				g.SetSummaryPokemonActive(g.opts.UserManager.GetUser().Party[g.Party.SelectedIndex])
				g.Party.SelectedIndex = -1
			}
			if g.Party.SummaryBtn.Hovered() {
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
						textStyle := material.Body2(g.Theme, "Summary")
						textStyle.Alignment = text.Middle
						textStyle.Color = TextColor

						return textStyle.Layout(gtx)
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if isBtnDisabled {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						}

						return g.Party.SummaryBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),

		// Move Up Btn
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			isBtnDisabled := g.Party.SelectedIndex == -1 || g.Party.SelectedIndex == 0

			btnBgColor := SecondaryBackgroundColor // Grey background for disabled
			if isBtnDisabled {
				btnBgColor = ButtonDisabledColor // Active background (blue)
			}

			if g.Party.MoveUpBtn.Clicked(gtx) {
				g.opts.UserManager.ChangePokemonOrder(g.opts.UserManager.GetUser().Party[g.Party.SelectedIndex], data.ChangeOrderMoveUp)
				g.Party.SelectedIndex = -1
			}
			if g.Party.MoveUpBtn.Hovered() {
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
						textStyle := material.Body2(g.Theme, "Move Up")
						textStyle.Alignment = text.Middle
						textStyle.Color = TextColor

						return textStyle.Layout(gtx)
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if isBtnDisabled {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						}

						return g.Party.MoveUpBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),

		// Move Down Btn
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			isBtnDisabled := g.Party.SelectedIndex == -1 || g.Party.SelectedIndex == len(g.opts.UserManager.GetUser().Party)-1

			btnBgColor := SecondaryBackgroundColor // Grey background for disabled
			if isBtnDisabled {
				btnBgColor = ButtonDisabledColor // Active background (blue)
			}

			if g.Party.MoveDownBtn.Clicked(gtx) {
				g.opts.UserManager.ChangePokemonOrder(g.opts.UserManager.GetUser().Party[g.Party.SelectedIndex], data.ChangeOrderMoveDown)
				g.Party.SelectedIndex = -1
			}
			if g.Party.MoveDownBtn.Hovered() {
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
						textStyle := material.Body2(g.Theme, "Move Down")
						textStyle.Alignment = text.Middle
						textStyle.Color = TextColor

						return textStyle.Layout(gtx)
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if isBtnDisabled {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						}

						return g.Party.MoveDownBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),
	)
}

func (g *Gui) renderPartyRowDetail(gtx layout.Context, pokemon *data.Pokemon) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		// Pokémon Image
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			img := loadImage(pokemon.SpritesURL.FrontPath)
			img.Fit = widget.Contain
			imgSize := image.Point{X: gtx.Dp(unit.Dp(64)), Y: gtx.Dp(unit.Dp(64))} // Fixed size
			return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max = imgSize
				return img.Layout(gtx)
			})
		}),
		// Name and HP bar
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset(layout.Inset{
				Top:    unit.Dp(8),
				Left:   unit.Dp(8),
				Right:  unit.Dp(16),
				Bottom: unit.Dp(8),
			}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name))
						name.Alignment = text.Start
						return name.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderHPBar(gtx, 1.0)
					}),
				)
			})
		}),
		// Pokémon Level
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			levelText := material.Body1(g.Theme, fmt.Sprintf("Lv. %v", pokemon.Level))
			levelText.Alignment = text.End
			return layout.Inset(layout.Inset{
				Top:    unit.Dp(8),
				Left:   unit.Dp(0),
				Right:  unit.Dp(64),
				Bottom: unit.Dp(8),
			}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return levelText.Layout(gtx)
			})
		}),
	)
}

func (g *Gui) renderEmptyPartyRow(gtx layout.Context) layout.Dimensions {
	// Render an empty row to maintain layout consistency
	rowHeight := gtx.Dp(unit.Dp(64)) // Same as Pokémon image height
	return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: rowHeight}}
}
