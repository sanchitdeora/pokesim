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
	"github.com/sanchitdeora/PokeSim/data"
)

func DefaultBoxProps(boxSize int) BoxProps {
	slog.Info("DefaultBoxProps", "boxSize", boxSize)

	pokmemonSelectionBtns := make([]*widget.Clickable, boxSize)
	for i := range pokmemonSelectionBtns {
		pokmemonSelectionBtns[i] = new(widget.Clickable)
	}
	return BoxProps{
		SelectedIndex:       -1,
		PokemonSelectedBtns: pokmemonSelectionBtns,
		SummaryBtn:          new(widget.Clickable),
		MoveToPartyBtn:      new(widget.Clickable),
	}
}

type BoxProps struct {
	SelectedIndex       int
	PokemonSelectedBtns []*widget.Clickable
	SummaryBtn          *widget.Clickable
	MoveToPartyBtn      *widget.Clickable
}

func (g *Gui) RenderBoxScreen(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,

		// Title
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Box")
			title.Font.Weight = font.Bold

			return title.Layout(gtx)
		}),
		layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						pokemons := g.opts.GameManager.GetBox()
						return g.renderBoxGallery(gtx, *pokemons)
					}))
			})
		}),

		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return g.renderBoxButtonRow(gtx)
		}),
	)
}

func (g *Gui) renderBoxGallery(gtx layout.Context, pokemons data.Box) layout.Dimensions {
	// Wrap the box pokemon list in a flex layout for better control

	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(g.Theme, g.BoxList).Layout(gtx, len(pokemons), func(gtx layout.Context, index int) layout.Dimensions {
				return g.renderBoxRow(gtx, pokemons, index)
			})
		}),
	)
}

func (g *Gui) renderBoxRow(gtx layout.Context, pokemons data.Box, index int) layout.Dimensions {
	// Create a row with up to 5 pokemons
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var children []layout.FlexChild
			for i := 0; i < 5; i++ {
				idx := index*5 + i
				if idx < len(pokemons) {
					pokemon := pokemons[idx].ToPokemon()
					children = append(children, layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {

						slog.Info("in loop", "index", index, "selectedIndex", g.Box.SelectedIndex, "gtx", gtx.Constraints)

						if g.Box.PokemonSelectedBtns[idx].Clicked(gtx) {
							slog.Info("pokemon Selected Clicked", "index", index, "selectedIndex", g.Box.SelectedIndex, "gtx", gtx.Constraints)
							g.Box.SelectedIndex = idx
						}
						return g.renderPokemonCard(gtx, pokemon, idx)
					}))
				} else { // Add placeholder for missing pokemons
					children = append(children, layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}))
				}
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderPokemonCard(gtx layout.Context, pokemon *data.Pokemon, index int) layout.Dimensions {
	if pokemon == nil {
		// Render an empty placeholder without a border
		return layout.Dimensions{}
	}

	isSelected := g.Box.SelectedIndex == index
	bgColor := SecondaryBackgroundColor // Default background
	if isSelected {
		bgColor = SelectedColor // Highlighted background
	}

	return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(250)))

		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
		fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, bgColor, 8)

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			// Pokemon Image
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				img := loadImage(pokemon.SpritesURL.FrontPath)
				img.Fit = widget.Contain
				return img.Layout(gtx)
			}),
			// Pokemon Name
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				name := material.Body1(g.Theme, pokemon.Name)
				name.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return name.Layout(gtx)
				})
			}),
			// Clickable Overlay for unlocked trainers
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return g.Box.PokemonSelectedBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: cardDims}
				})

			}),
		)
	})
}

func (g *Gui) renderBoxButtonRow(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),

		// Summary Btn
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			isBtnDisabled := g.Box.SelectedIndex == -1

			btnBgColor := SecondaryBackgroundColor // Grey background for disabled
			if isBtnDisabled {
				btnBgColor = ButtonDisabledColor // Active background (blue)
			}

			if g.Box.SummaryBtn.Clicked(gtx) {
				g.SetSummaryPokemonActive(g.getSelectedBoxPokemonByIndex(g.Box.SelectedIndex))
				g.Box.SelectedIndex = -1
			}
			if g.Box.SummaryBtn.Hovered() {
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

						return g.Box.SummaryBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),

		// Move To Party Btn
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			isBtnDisabled := g.Box.SelectedIndex == -1

			btnBgColor := SecondaryBackgroundColor // Grey background for disabled
			if isBtnDisabled {
				btnBgColor = ButtonDisabledColor // Active background (blue)
			}

			if g.Box.MoveToPartyBtn.Clicked(gtx) {
				g.opts.BoxManager.Swap(g.getSelectedBoxPokemonByIndex(g.Box.SelectedIndex), nil)
				g.Box = DefaultBoxProps(len(*g.opts.GameManager.GetBox()))
			}
			if g.Box.MoveToPartyBtn.Hovered() {
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
						textStyle := material.Body2(g.Theme, "Move To Party")
						textStyle.Alignment = text.Middle
						textStyle.Color = TextColor

						return textStyle.Layout(gtx)
					}),
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if isBtnDisabled {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						}

						return g.Box.MoveToPartyBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),
	)
}

func (g *Gui) getSelectedBoxPokemonByIndex(index int) *data.Pokemon {
	return (*g.opts.GameManager.GetBox())[index].ToPokemon()
}
