package gui

import (
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func (g *Gui) RenderBagScreen(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,

		// Title
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Bag")
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
						return g.renderBagList(gtx)
					}))
			})
		}),
	)
}

func (g *Gui) renderBagList(gtx layout.Context) layout.Dimensions {
	// Get the user's bag data
	bag := g.opts.UserManager.GetUser().Bag

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var bagItems []layout.FlexChild
			for name, item := range bag {
				// Render items in bag
				bagItems = append(bagItems, layout.Flexed(1.0/7.0, func(gtx layout.Context) layout.Dimensions {
					return g.renderBagRow(gtx, name, &item)
				}))
			}
			// slog.Info("bagItems", "len", len(bagItems), "item", bagItems, "bag", bag)
			return bagItems
		}(gtx)...,
	)
}

func (g *Gui) renderBagRow(gtx layout.Context, name data.ItemName, item *data.Item) layout.Dimensions {
	if item == nil {
		// Render an empty placeholder without a border
		return layout.Dimensions{}
	}

	// Render a single row for a item in the bag
	return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := gtx.Constraints.Max

		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
		fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, SecondaryBackgroundColor, 8)

		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			// Item Image
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				img := loadImage("./assets/utils/pokeball.png")
				img.Fit = widget.Contain
				img.Position = layout.Center
				imgSize := image.Point{X: gtx.Dp(unit.Dp(64)), Y: gtx.Dp(unit.Dp(64))} // Fixed size
				return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max = imgSize
					return img.Layout(gtx)
				})
			}),
			// Name
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset(layout.Inset{
					Top:    unit.Dp(0),
					Left:   unit.Dp(0),
					Right:  unit.Dp(0),
					Bottom: unit.Dp(0),
				}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(string(name)))
							name.Alignment = text.Start
							return name.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							name := material.Caption(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(item.Description))
							name.Alignment = text.Start
							return name.Layout(gtx)
						}),
					)
				})
			}),
			// Pokémon Level
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				levelText := material.Body1(g.Theme, fmt.Sprintf("x%v", item.Count))
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
	})
}
