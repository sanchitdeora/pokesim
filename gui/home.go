package gui

import (
	"fmt"
	"image"
	"log/slog"

	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (g *Gui) RenderHomeScreen(gtx layout.Context) layout.Dimensions {
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
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {

				title := material.H4(g.Theme, "Home")
				title.Font.Weight = font.Bold

				return title.Layout(gtx)
			}),

			layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {

				return layout.Inset(layout.Inset{
					Top: unit.Dp(20),
				}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderStatCard(gtx)
				})
			}),

			layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset(layout.Inset{
							Top: unit.Dp(10),
						}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return g.renderActionCard(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset(layout.Inset{
							Top: unit.Dp(15),
						}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return g.renderQuickAccess(gtx)
						})
					}),
				)
			}),
		)
	})
}

func (g *Gui) renderStatCard(gtx layout.Context) layout.Dimensions {
	stats := []struct {
		mainText string
		label    string
	}{
		{g.opts.UserManager.GetUser().Name, "Name"},
		{strconv.Itoa(len(g.opts.UserManager.GetUser().Stats.Badges)), "Badges"},
		{strconv.Itoa(g.opts.UserManager.GetUser().Money), "Money"},
	}

	cardPadding := gtx.Dp(unit.Dp(8))

	// Create a horizontal layout for the three boxes
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween, // Add even spacing between boxes
	}.Layout(gtx,
		// Create a FlexChild for each stat box
		func(gtx layout.Context) []layout.FlexChild {
			children := []layout.FlexChild{}

			for idx, stat := range stats {
				// Add a Flexed box that evenly splits the space

				var leftPadding, rightPadding unit.Dp
				if idx > 0 {
					leftPadding = unit.Dp(cardPadding)
				}
				if idx < len(stats)-1 {
					rightPadding = unit.Dp(cardPadding)
				}

				children = append(children, layout.Flexed(1.0/float32(len(stats)), func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Left:  leftPadding,
						Right: rightPadding,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// Box size and border properties
						borderSize := image.Pt(gtx.Constraints.Max.X, gtx.Dp(100))
						drawRoundedBorder(gtx, borderSize, CardBorderColor, unit.Dp(1), unit.Dp(8))

						// Layout with vertical alignment
						return layout.Flex{
							Axis:      layout.Vertical,
							Alignment: layout.Start,
						}.Layout(gtx,
							// Spacer to push the main text to the vertical center
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									textStyle := material.H5(g.Theme, stat.mainText)
									textStyle.Alignment = text.Middle
									textStyle.Font.Weight = font.Bold
									return textStyle.Layout(gtx)
								})
							}),
							// Caption text at the bottom
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								labelStyle := material.Caption(g.Theme, stat.label)
								labelStyle.Alignment = text.Middle
								labelStyle.Color = TextColor // Gray for label
								return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, labelStyle.Layout)
							}),
						)
					})
				}))
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderActionCard(gtx layout.Context) layout.Dimensions {
	actions := []struct {
		name   string
		icon   *widget.Icon // Placeholder for icon (if needed)
		screen Screen       // Navigation target
	}{
		{"Trainers", nil, TrainerScreen},
		{"Wild", nil, WildScreen},
		{"Tournaments", nil, ToBeReplaced},
		{"PokéShop", nil, ToBeReplaced},
		{"Save", nil, ToBeReplaced},
	}

	cardHeight := gtx.Dp(unit.Dp(80))
	cardPadding := gtx.Dp(unit.Dp(8)) // Gap between cards
	rowMargin := gtx.Dp(unit.Dp(10))  // Space between rows
	cardsPerRow := 5                  // Number of cards per row

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		totalWidth := gtx.Constraints.Max.X // Available width in the container

		// Calculate the total gap space between cards and adjust for padding
		totalGapSpace := 2 * (cardPadding * (cardsPerRow - 1))
		cardWidth := (totalWidth - totalGapSpace) / cardsPerRow

		// Create a vertical layout for each row of cards
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEnd,
		}.Layout(gtx,
			func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				for i := 0; i < len(actions); i += cardsPerRow {
					rowItems := actions[i:min(i+cardsPerRow, len(actions))]

					children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Create space between rows
						return layout.Inset{
							Bottom: unit.Dp(rowMargin),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceBetween, // Space between cards
							}.Layout(gtx,
								func(gtx layout.Context) []layout.FlexChild {
									var rowChildren []layout.FlexChild

									for idx, action := range rowItems {
										action := action // Properly scope the loop variable

										// For the first and last item, no extra padding is added, only the internal space
										var leftPadding, rightPadding unit.Dp
										if idx > 0 {
											leftPadding = unit.Dp(cardPadding)
										}
										if idx < len(rowItems)-1 {
											rightPadding = unit.Dp(cardPadding)
										}

										// Apply the calculated card width
										rowChildren = append(rowChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{
												Left:   leftPadding,          // Space before the card (between cards)
												Right:  rightPadding,         // Space after the card (between cards)
												Top:    unit.Dp(cardPadding), // Optional top padding
												Bottom: unit.Dp(cardPadding), // Optional bottom padding
											}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												// Define constraints for the card, using the adjusted card width
												cardGtx := gtx
												cardGtx.Constraints = layout.Constraints{
													Min: image.Pt(cardWidth, cardHeight),
													Max: image.Pt(cardWidth, cardHeight),
												}
												drawRoundedBorder(cardGtx, image.Pt(cardWidth, cardHeight), CardBorderColor, unit.Dp(1), unit.Dp(8))

												// Layout the card content
												return layout.Center.Layout(cardGtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Stack{
														Alignment: layout.Center,
													}.Layout(gtx,
														layout.Stacked(func(gtx layout.Context) layout.Dimensions {
															textStyle := material.Body1(g.Theme, action.name)
															textStyle.Alignment = text.Middle
															textStyle.Font.Weight = font.Bold
															return textStyle.Layout(gtx)
														}),
														layout.Stacked(func(gtx layout.Context) layout.Dimensions {
															if action.icon != nil {
																return action.icon.Layout(gtx, TextColor) // Example icon size
															}
															return layout.Dimensions{}
														}),
														// Clickable Overlay for unlocked trainers
														layout.Expanded(func(gtx layout.Context) layout.Dimensions {
															if g.Buttons[action.screen].Clicked(gtx) {
																slog.Info("Clicked", "page", action.screen)
																g.SetCurrentScreen(action.screen)
															}
															return g.Buttons[action.screen].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
																return layout.Dimensions{Size: gtx.Constraints.Max}
															})
														}),
													)
												})
											})
										}))
									}
									return rowChildren
								}(gtx)...,
							)
						})
					}))
				}
				return children
			}(gtx)...,
		)
	})
}

func (g *Gui) renderQuickAccess(gtx layout.Context) layout.Dimensions {
	btnWidth := gtx.Dp(unit.Dp(100)) // Fixed width for the button

	i1, _ := widget.NewIcon([]byte{})
	i2, _ := widget.NewIcon([]byte{})
	i3, _ := widget.NewIcon([]byte{})

	lenParty := len(g.opts.UserManager.GetUser().Party)
	itemCount := g.opts.UserManager.GetUser().GetBagItemCount()

	quickAccessItems := []struct {
		icon      *widget.Icon
		title     string
		subtext   string
		buttonTxt string
		page      Screen
	}{
		{icon: i1, title: "Party", subtext: fmt.Sprintf("You have %v/6 Pokemon", lenParty), buttonTxt: "View Party", page: PartyScreen},
		{icon: i2, title: "Bag", subtext: fmt.Sprintf("You have %v items", itemCount), buttonTxt: "View Bag", page: Bag},
		{icon: i3, title: "Pokédex", subtext: fmt.Sprintf("%v Pokemon seen", 7), buttonTxt: "Open", page: ToBeReplaced},
	}

	// slog.Info("Total Width", "gtx", gtx.Constraints.Max, "gtx min", gtx.Constraints.Min)

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceStart,
		}.Layout(gtx,
			func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					textStyle := material.H5(g.Theme, "Quick Access")
					textStyle.Alignment = text.Middle
					textStyle.Font.Weight = font.Bold
					return layout.Inset(layout.Inset{
						Bottom: unit.Dp(30),
					}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return textStyle.Layout(gtx)
					})
				}))

				for idx, item := range quickAccessItems {
					// Apply inset between items except for the first one
					children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {

						topPadding := unit.Dp(0)
						if idx > 0 {
							topPadding = unit.Dp(8)
						}
						return layout.Inset{
							Top: topPadding, // Adjust the value for the inset
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceBetween,
							}.Layout(gtx,
								func(gtx layout.Context) []layout.FlexChild {
									var rowChildren []layout.FlexChild

									rowChildren = append(rowChildren, layout.Flexed(0.85, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Horizontal,
										}.Layout(gtx,
											func(gtx layout.Context) []layout.FlexChild {
												var children []layout.FlexChild

												children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{
														Axis: layout.Vertical,
													}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															textStyle := material.Body1(g.Theme, item.title)
															return textStyle.Layout(gtx)
														}),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															textStyle := material.Body2(g.Theme, item.subtext)
															return textStyle.Layout(gtx)
														}),
													)
												}))

												return children
											}(gtx)...,
										)
									}))

									// Right button (extreme right)
									rowChildren = append(rowChildren, layout.Flexed(0.15, func(gtx layout.Context) layout.Dimensions {
										btnStyle := material.Button(g.Theme, g.Buttons[item.page], item.buttonTxt)
										if g.Buttons[item.page].Clicked(gtx) {
											slog.Info("Clicked", "page", item.page)
											g.SetCurrentScreen(item.page)
										}
										btnStyle.Color = TextColor
										btnStyle.Background = SecondaryBackgroundColor

										btnLayout := btnStyle.Layout(gtx)
										btnLayout.Size.X = btnWidth // Set fixed width for the button

										return btnLayout
									}))

									return rowChildren
								}(gtx)...,
							)
						})
					}))
				}

				return children
			}(gtx)...,
		)
	})
}
