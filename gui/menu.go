package gui

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (g *Gui) RenderMenu(gtx layout.Context) layout.Dimensions {
	// Set the constraints for the menu
	// gtx.Constraints.Min = image.Pt(gtx.Dp(unit.Dp(1300)), gtx.Dp(unit.Dp(0)))
	// gtx.Constraints.Max = image.Pt(gtx.Dp(unit.Dp(1300)), gtx.Dp(unit.Dp(800)))

	insets := layout.UniformInset(unit.Dp(1))

	return insets.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Horizontal,
			// Spacing:   layout.SpaceBetween,
			Alignment: layout.Start, // Add space between menu items
		}.Layout(gtx,
			// left side of the menu - App Name
			g.renderAppName(),

			// right side of the menu with navbar
			g.renderMenuItems(),
		)
	})
}

func (g *Gui) renderAppName() layout.FlexChild {
	return layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
		titleStyle := material.H6(g.Theme, ":) PokéSim")
		titleStyle.Font.Weight = font.Bold
		titleStyle.Color = TextColor

		return layout.Inset(layout.Inset{
			Top:  unit.Dp(10),
			Left: unit.Dp(10),
		}).Layout(gtx, titleStyle.Layout)
	})
}

func (g *Gui) renderMenuItems() layout.FlexChild {
	menuItems := []struct {
		label  string
		screen Screen
	}{
		{"Home", HomeScreen},
		{"Trainer", TrainerScreen},
		{"Bag", BagScreen},
		{"Box", BoxScreen},
	}

	return layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.End,
		}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
			children := []layout.FlexChild{}

			// Add menu items
			for _, item := range menuItems {
				item := item // avoid capturing loop variable
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := g.Buttons[item.screen]

					if btn.Clicked(gtx) {
						if g.isMenuActive() {
							g.SetCurrentScreen(item.screen)
						} else {
							// do nothing
						}
					}

					return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Label(g.Theme, unit.Sp(15), item.label)
						label.Color = TextColor

						return layout.UniformInset(unit.Dp(15)).Layout(gtx, label.Layout)
					})
				}))
			}
			return children
		}(gtx)...)
	})
}

func (g *Gui) RenderCurrentScreen(gtx layout.Context) layout.Dimensions {
	// slog.Info("Rendering current window", "current window", a.CurrentPage)
	switch g.CurrentWindow {
	case NewGameScreen:
		return g.RenderNewGameScreen(gtx)
	case LoadGameScreen:
		return g.RenderLoadGameScreen(gtx)
	case HomeScreen:
		return g.RenderHomeScreen(gtx)
	case BattleScreen:
		return g.RenderBattleScreen(gtx)
	case TrainerScreen:
		return g.RenderTrainerScreen(gtx)
	case WildScreen:
		return g.RenderWildScreen(gtx)
	case PartyScreen:
		return g.RenderPartyScreen(gtx)
	case PokemonSummaryScreen:
		return g.RenderPokemonSummary(gtx)
	case BagScreen:
		return g.RenderBagScreen(gtx)
	case BoxScreen:
		return g.RenderBoxScreen(gtx)
	case ShopScreen:
		return g.RenderShopScreen(gtx)
	case ToBeReplaced:
		return g.RenderBoxScreen(gtx)
	default:
		return layout.Dimensions{}
	}
}

func (g *Gui) SetCurrentScreen(newWindow Screen) {
	g.CurrentWindow = newWindow
}

func (g *Gui) isMenuActive() bool {
	return g.opts.GameManager.Get() != nil
}
