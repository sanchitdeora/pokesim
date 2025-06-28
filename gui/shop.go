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

type ShopBtns struct {
	decreaseBtn *widget.Clickable
	quantity    int
	increaseBtn *widget.Clickable
	purchaseBtn *widget.Clickable
	sellBtn     *widget.Clickable
}

type ShopProps struct {
	btns map[data.ItemName]ShopBtns
}

func DefaultShopProps() ShopProps {
	var shopProps ShopProps

	shopItems := getItemShopList()
	shopProps.btns = make(map[data.ItemName]ShopBtns, len(shopItems))
	for _, item := range shopItems {
		shopProps.btns[item.StoreItem.Name] = ShopBtns{
			decreaseBtn: new(widget.Clickable),
			quantity:    item.StoreItem.Item.Count,
			increaseBtn: new(widget.Clickable),
			purchaseBtn: new(widget.Clickable),
			sellBtn:     new(widget.Clickable),
		}
	}

	return shopProps
}

func (g *Gui) RenderShopScreen(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Title
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "PokéShop")
			title.Font.Weight = font.Bold

			return title.Layout(gtx)
		}),
		layout.Flexed(0.9, func(gtx layout.Context) layout.Dimensions {
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
						return g.renderShopItemList(gtx)
					}))
			})
		}),
	)
}

func (g *Gui) renderShopItemList(gtx layout.Context) layout.Dimensions {
	items := getItemShopList()
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var rows []layout.FlexChild
			for i, item := range items {
				index := i
				// Render the Pokémon in the party
				// slog.Info("in loop", "index", index, "selectedIndex", g.Party.SelectedIndex, "gtx", gtx.Constraints)
				rows = append(rows, layout.Flexed(1.0/6.0, func(gtx layout.Context) layout.Dimensions {
					if g.Party.PokemonSelectedBtns[index].Clicked(gtx) {
						g.Party.SelectedIndex = index
						// slog.Info("pokemon Selected Clicked", "index", index, "selectedIndex", g.Party.SelectedIndex, "gtx", gtx.Constraints)
					}
					return g.renderItemRow(gtx, item)
				}))

			}
			return rows
		}(gtx)...,
	)
}

func (g *Gui) renderItemRow(gtx layout.Context, item ItemShopUI) layout.Dimensions {
	// Render a single row for an item in shop
	return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := gtx.Constraints.Max

		drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
		fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, SecondaryBackgroundColor, 8)

		return g.renderItemRowDetail(gtx, item)
	})
}

func (g *Gui) renderItemRowDetail(gtx layout.Context, item ItemShopUI) layout.Dimensions {
	shopBtns, ok := g.Shop.btns[item.StoreItem.Name]
	if !ok {
		return layout.Dimensions{}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Item Image
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			img := loadImage(item.ImageURL)
			img.Fit = widget.Contain
			imgSize := image.Point{X: gtx.Dp(unit.Dp(48)), Y: gtx.Dp(unit.Dp(48))} // Fixed size
			return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max = imgSize
				return img.Layout(gtx)
			})
		}),
		// Name and Description
		layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset(layout.Inset{
				Top:    unit.Dp(8),
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
				Bottom: unit.Dp(8),
			}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(string(item.StoreItem.Name)))
						name.Alignment = text.Start
						return name.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						name := material.Body2(g.Theme, string(item.StoreItem.Item.Description))
						name.Alignment = text.Start
						return name.Layout(gtx)
					}),
				)
			})
		}),
		// Quantity, Cost
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceAround,
			}.Layout(gtx,
				// Quantity Selector
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return g.renderQuantitySelector(gtx, &shopBtns, item.StoreItem.Name) // Custom quantity selector
					})
				}),

				// Item Cost
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					cost := material.Body2(g.Theme, fmt.Sprintf("%d ₽", (item.StoreItem.Item.CostPrice)*shopBtns.quantity))
					cost.Color = TextColor
					cost.Alignment = text.Middle
					return layout.Center.Layout(gtx, cost.Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(g.Theme, shopBtns.purchaseBtn, "Buy")
				if shopBtns.purchaseBtn.Clicked(gtx) {
					slog.Info("buy button clicked", "item", item.StoreItem.Name)
					g.PurchaseItem(item)
					shopBtns.quantity = 1
					g.Shop.btns[item.StoreItem.Name] = shopBtns
				}
				return btn.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(g.Theme, shopBtns.sellBtn, "Sell")
				btn.Background = RedBtnColor
				if shopBtns.sellBtn.Clicked(gtx) {
					slog.Info("sell button clicked", "item", item.StoreItem.Name)
					g.SellItem(item)
					shopBtns.quantity = 1
					g.Shop.btns[item.StoreItem.Name] = shopBtns
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

// renderQuantitySelector for increment/decrement
func (g *Gui) renderQuantitySelector(gtx layout.Context, shopBtns *ShopBtns, name data.ItemName) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(g.Theme, shopBtns.decreaseBtn, "-")
			if shopBtns.decreaseBtn.Clicked(gtx) && shopBtns.quantity > 0 {
				shopBtns.quantity--
				g.Shop.btns[name] = *shopBtns
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			txt := material.Body1(g.Theme, fmt.Sprintf("%d", shopBtns.quantity))
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return layout.Center.Layout(gtx, txt.Layout) })
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(g.Theme, shopBtns.increaseBtn, "+")
			if shopBtns.increaseBtn.Clicked(gtx) {
				shopBtns.quantity++
				g.Shop.btns[name] = *shopBtns
			}
			return btn.Layout(gtx)
		}),
	)
}

func (g *Gui) PurchaseItem(item ItemShopUI) {
	item.StoreItem.Item.Count = g.Shop.btns[item.StoreItem.Name].quantity
	err := g.opts.UserManager.PurchaseItem(item.StoreItem.Item)
	if err != nil {
		slog.Error("failed to purchase item", "err", err)
	}
}

func (g *Gui) SellItem(item ItemShopUI) {
	item.StoreItem.Item.Count = g.Shop.btns[item.StoreItem.Name].quantity
	err := g.opts.UserManager.SellItem(item.StoreItem.Item)
	if err != nil {
		slog.Error("failed to sell item", "err", err)
	}
}
