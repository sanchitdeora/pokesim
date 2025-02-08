package gui

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func (g *Gui) RenderBattleScreen(gtx layout.Context) layout.Dimensions {
	// Use a vertical layout for screen elements
	xInset := gtx.Dp(unit.Dp(200))
	yInset := gtx.Dp(unit.Dp(30))

	mainScreen := layout.Inset(layout.Inset{
		Top:    unit.Dp(yInset),
		Left:   unit.Dp(xInset),
		Right:  unit.Dp(xInset),
		Bottom: unit.Dp(yInset),
	}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					// Top Row: Opponent Info
					layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
						return g.renderOpponentSection(gtx)
					}),
					// Middle Row: User Info
					layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
						return g.renderUserSection(gtx)
					}),
					// Bottom Section: Battle Options
					layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
						return g.renderBattleOptions(gtx)
					}),
				)
			}),
		)
	})

	switch g.TrainerBattle.DialogActionArea {
	case SwitchDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset(layout.Inset{
					Top:    unit.Dp(yInset),
					Left:   unit.Dp(xInset),
					Right:  unit.Dp(xInset),
					Bottom: unit.Dp(yInset),
				}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return g.renderSwitchDialog(gtx)
						}),
					)
				})
			}),
		)
	case BagDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset(layout.Inset{
					Top:    unit.Dp(yInset),
					Left:   unit.Dp(xInset),
					Right:  unit.Dp(xInset),
					Bottom: unit.Dp(yInset),
				}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return g.renderBagDialog(gtx)
						}),
					)
				})
			}),
		)
	default:
		return mainScreen
	}
}

func (g *Gui) renderOpponentSection(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// Opponent Info
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			borderSize := image.Pt(gtx.Constraints.Max.X, gtx.Dp(80))
			drawRoundedBorder(gtx, borderSize, CardBorderColor, unit.Dp(1), unit.Dp(8))

			return g.renderOpponentInfo(gtx)
		}),
		// Opponent Image
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {

			return g.renderOpponentImage(gtx)
		}),
	)
}

// Opponent Info
func (g *Gui) renderOpponentInfo(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// Name and Level
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(g.TrainerBattle.Opponent.GetActivePokemon().Name))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Lv. %v", g.TrainerBattle.Opponent.GetActivePokemon().Level))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// HP Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderHPBar(gtx, g.GetRemainingHPFraction(g.TrainerBattle.Opponent.GetActivePokemon())) // Example: 50% HP remaining
				})
			}),
		)
	})
}

// Opponent Image
func (g *Gui) renderOpponentImage(gtx layout.Context) layout.Dimensions {
	img := g.loadPokemonImage(g.TrainerBattle.Opponent.GetActivePokemon().SpritesURL.FrontPath) // Replace with your image loading logic
	img.Fit = widget.Contain
	img.Position = layout.Center
	return img.Layout(gtx)
}

func (g *Gui) renderUserSection(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// User Image
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return g.renderUserImage(gtx)
		}),
		// User Info
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			borderSize := image.Pt(gtx.Constraints.Max.X, gtx.Dp(140))
			drawRoundedBorder(gtx, borderSize, CardBorderColor, unit.Dp(1), unit.Dp(8))
			return g.renderUserInfo(gtx)
		}),
	)
}

// User Info
func (g *Gui) renderUserInfo(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// Name and Level
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(g.TrainerBattle.User.GetActivePokemon().Name))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Lv. %v", g.TrainerBattle.User.GetActivePokemon().Level))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// HP Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderHPBar(gtx, g.GetRemainingHPFraction(g.TrainerBattle.User.GetActivePokemon())) // Example: 50% HP remaining
				})
			}),
			// Remaining HP
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, "")
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)

					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(g.Theme, fmt.Sprintf("%d/%d HP", g.TrainerBattle.User.GetActivePokemon().BattleHP, g.TrainerBattle.User.GetActivePokemon().Stats.HP.Value))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// Experience Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderEXPBar(gtx, g.GetRemainingEXPFraction(g.TrainerBattle.User.GetActivePokemon()))
				})
			}),
		)
	})
}

// User Image
func (g *Gui) renderUserImage(gtx layout.Context) layout.Dimensions {
	img := g.loadPokemonImage(g.TrainerBattle.User.GetActivePokemon().SpritesURL.BackPath) // Replace with your image loading logic
	img.Fit = widget.Contain
	img.Position = layout.Center
	return img.Layout(gtx)
}

// Battle Options
func (g *Gui) renderBattleOptions(gtx layout.Context) layout.Dimensions {

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(0.6,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset(layout.Inset{Right: unit.Dp(15)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					drawRoundedBorder(gtx, image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y), CardBorderColor, unit.Dp(1), unit.Dp(8))

					return g.renderBattleLog(gtx)
				})
			}),
		layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
			drawRoundedBorder(gtx, image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y), CardBorderColor, unit.Dp(1), unit.Dp(8))

			return g.renderActionBtns(gtx)
		}),
	)
}

func (g *Gui) renderBattleAction(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
		var children []layout.FlexChild

		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[AttackBtn].Clicked(gtx) {
							g.SetActionBtns(Attack)
						}
						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[AttackBtn], "Attack")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[SwitchBtn].Clicked(gtx) {
							g.TrainerBattle.DialogActionArea = SwitchDialog
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[SwitchBtn], "Switch")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))

				return children
			}(gtx)...,
			)
		}))
		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[BagBtn].Clicked(gtx) {
							g.TrainerBattle.DialogActionArea = BagDialog
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[BagBtn], "Bag")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[RunBtn].Clicked(gtx) {
							g.TrainerBattle.LogChan <- "No! There's no running from a Trainer Battle!"
							g.SetActionBtns(Actions)
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[RunBtn], "Run")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))

				return children
			}(gtx)...,
			)
		}))

		return children
	}(gtx)...,
	)
}

// when clicking attack, should display this. Will add this later
func (g *Gui) renderBattleAttack(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
		var children []layout.FlexChild

		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[Move1Btn].Clicked(gtx) {
							g.performAttack(g.TrainerBattle.User.GetActivePokemon().Moveset.Move1)
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[Move1Btn], "Fire Punch")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[Move2Btn].Clicked(gtx) {
							g.performAttack(g.TrainerBattle.User.GetActivePokemon().Moveset.Move2)
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[Move2Btn], "Ice Punch")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))

				return children
			}(gtx)...,
			)
		}))
		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[Move3Btn].Clicked(gtx) {
							g.performAttack(g.TrainerBattle.User.GetActivePokemon().Moveset.Move3)
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[Move3Btn], "Thunder Punch")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.TrainerBattle.ActionButtons[Move4Btn].Clicked(gtx) {
							g.performAttack(g.TrainerBattle.User.GetActivePokemon().Moveset.Move4)
						}

						btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[Move4Btn], "Mach Punch")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))

				return children
			}(gtx)...,
			)
		}))
		return children
	}(gtx)...,
	)
}

// Switch Dialog
func (g *Gui) renderSwitchDialog(gtx layout.Context) layout.Dimensions {

	drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))
	fillRoundedShape(gtx, image.Rectangle{Max: gtx.Constraints.Max}, SecondaryBackgroundColor, 8)

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Overlay background
			col := color.NRGBA{R: 0, G: 0, B: 0, A: 10} // Semi-transparent black
			return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.Fill(gtx.Ops, col)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Dialog box
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{Alignment: layout.Center}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceEvenly,
						}.Layout(gtx,
							// Title
							layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
								label := material.H6(g.Theme, "Switch Pokémon")
								label.Alignment = text.Middle
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							// Pokémon list
							layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {

								return material.List(g.Theme, &widget.List{
									List: layout.List{Axis: layout.Vertical},
								}).Layout(gtx, len(g.TrainerBattle.User.GetParty()), func(gtx layout.Context, i int) layout.Dimensions {
									pokemon := g.TrainerBattle.User.GetParty()[i]

									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										// Pokémon Image
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											img := g.loadPokemonImage(pokemon.SpritesURL.FrontPath)
											img.Fit = widget.Contain
											imgSize := image.Point{X: gtx.Dp(unit.Dp(64)), Y: gtx.Dp(unit.Dp(64))} // Fixed size
											return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
														name := material.Body1(g.Theme,
															fmt.Sprintf("%s (Lv. %d)", utils.ToCapitalizeFirstLetterOfEachWord(pokemon.Name), pokemon.Level))
														name.Alignment = text.Start
														return name.Layout(gtx)
													}),
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														return g.renderHPBar(gtx, float32(pokemon.BattleHP)/float32(pokemon.Stats.HP.Value))
													}),
												)
											})
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return layout.Inset(layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												btnStyle := material.Button(g.Theme, g.TrainerBattle.PokemonSwitchButtons[i], "Switch")
												if g.TrainerBattle.PokemonSwitchButtons[i].Clicked(gtx) {
													slog.Info("Switching to Pokémon", "pokemon", pokemon.Name)
													if g.TrainerBattle.User.GetParty()[i].BattleHP > 0 {
														g.performSwitch(i)
														g.TrainerBattle.DialogActionArea = MainBattle
													} else {
														slog.Info("Can't switch to fainted Pokémon")
													}
												}

												btnStyle.Color = TextColor
												btnStyle.Background = PrimaryBackgroundColor

												btnLayout := btnStyle.Layout(gtx)
												btnLayout.Size.X = gtx.Dp(unit.Dp(128)) // Set fixed width for the button

												return btnLayout
											})
										}),
									)
								})
							}),

							// Cancel button
							layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(g.Theme, g.TrainerBattle.ActionButtons[CancelSwitchBtn], "Cancel")
								btn.Background = color.NRGBA{R: 200, G: 0, B: 0, A: 255} // Red for cancel
								if g.TrainerBattle.ActionButtons[CancelSwitchBtn].Clicked(gtx) {
									g.TrainerBattle.DialogActionArea = MainBattle
								}
								return layout.Center.Layout(gtx, btn.Layout)
							}),
						)
					}),
				)
			})
		}),
	)
}

// Bag Dialog
func (g *Gui) renderBagDialog(gtx layout.Context) layout.Dimensions {
	drawRoundedBorder(gtx, gtx.Constraints.Max, CardBorderColor, unit.Dp(1), unit.Dp(8))
	fillRoundedShape(gtx, image.Rectangle{Max: gtx.Constraints.Max}, SecondaryBackgroundColor, 8)

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Dialog box
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{Alignment: layout.Center}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceEvenly,
						}.Layout(gtx,
							// Title
							layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
								label := material.H6(g.Theme, "Bag")
								label.Alignment = text.Middle
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							// Item list
							layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {

								return material.List(g.Theme, &widget.List{
									List: layout.List{Axis: layout.Vertical},
								}).Layout(gtx, len(data.AllItems), func(gtx layout.Context, i int) layout.Dimensions {
									itemName := data.AllItems[i]
									item, exists := g.TrainerBattle.User.GetTrainer().Bag[itemName]
									if !exists {
										return layout.Dimensions{}
									}

									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										// // Item Image
										// layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
										// 	img := g.loadPokemonImage(pokemon.SpritesURL.FrontPath)
										// 	img.Fit = widget.Contain
										// 	imgSize := image.Point{X: gtx.Dp(unit.Dp(64)), Y: gtx.Dp(unit.Dp(64))} // Fixed size
										// 	return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										// 		gtx.Constraints.Max = imgSize
										// 		return img.Layout(gtx)
										// 	})
										// }),
										// Name and Description
										layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
											slog.Info("Item inside 1", "name", itemName, "count", item.Count, "gtx", gtx.Constraints.Max)

											return layout.Inset(layout.Inset{
												Top:    unit.Dp(8),
												Left:   unit.Dp(8),
												Right:  unit.Dp(16),
												Bottom: unit.Dp(8),
											}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(string(itemName)))
														name.Alignment = text.Start
														return name.Layout(gtx)
													}),
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														name := material.Body2(g.Theme, item.Description)
														name.Alignment = text.Start
														return name.Layout(gtx)
													}),
												)
											})
										}),
										layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
											slog.Info("Item inside 2", "name", itemName, "count", item.Count, "gtx", gtx.Constraints.Max)

											label := material.Body1(g.Theme, fmt.Sprintf("x%v", item.Count))
											label.Alignment = text.Middle
											// label.Font.Weight = font.Bold
											return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return label.Layout(gtx)
											})

										}),
										layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
											slog.Info("Item inside 3", "name", itemName, "count", item.Count, "gtx", gtx.Constraints.Max)

											return layout.Inset(layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												btnStyle := material.Button(g.Theme, g.TrainerBattle.UseItemButtons[itemName], "Use")
												if g.TrainerBattle.UseItemButtons[itemName].Clicked(gtx) {
													slog.Info("Using item", "item", itemName)
													if g.TrainerBattle.User.GetTrainer().Bag[itemName].Category == data.MedicalItems && (g.TrainerBattle.User.GetActivePokemon().BattleHP < g.TrainerBattle.User.GetActivePokemon().Pokemon.Stats.HP.Value) {
														g.performUseItem(&item)
														g.TrainerBattle.DialogActionArea = MainBattle
													}
												}

												btnStyle.Color = TextColor
												btnStyle.Background = PrimaryBackgroundColor

												btnLayout := btnStyle.Layout(gtx)
												btnLayout.Size.X = gtx.Dp(unit.Dp(128)) // Set fixed width for the button

												return btnLayout
											})
										}),
									)
								})
							}),

							// Cancel button
							layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(g.Theme, g.TrainerBattle.ActionButtons[CancelSwitchBtn], "Cancel")
								btn.Background = color.NRGBA{R: 200, G: 0, B: 0, A: 255} // Red for cancel
								if g.TrainerBattle.ActionButtons[CancelSwitchBtn].Clicked(gtx) {
									g.TrainerBattle.DialogActionArea = MainBattle
								}
								return layout.Center.Layout(gtx, btn.Layout)
							}),
						)
					}),
				)
			})
		}),
	)
}

func (g *Gui) renderBattleEnd(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if g.TrainerBattle.ActionButtons[EndBattleBtn].Clicked(gtx) {
			g.SetCurrentScreen(HomeScreen)
		}

		btnStyle := material.Button(g.Theme, g.TrainerBattle.ActionButtons[EndBattleBtn], "End Battle")
		btnStyle.Background = SecondaryBackgroundColor
		btnStyle.Color = TextColor

		return btnStyle.Layout(gtx)
	})
}

func (g *Gui) renderBattleLog(gtx layout.Context) layout.Dimensions {

	// Wrap the trainer list in a flex layout for better control
	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Set the available height for the list
			availableHeight := gtx.Constraints.Max.Y

			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.List(g.Theme, g.TrainerBattle.LogList).Layout(gtx, len(g.TrainerBattle.LogContent), func(gtx layout.Context, index int) layout.Dimensions {
					gtx.Constraints.Max.Y = availableHeight
					return material.Body2(g.Theme, g.TrainerBattle.LogContent[index]).Layout(gtx)
				})
			})
		}),
	)
}

func (g *Gui) performAttack(move *data.Moves) {
	g.TrainerBattle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: g.TrainerBattle.User.GetActivePokemon(),
		Target:   g.TrainerBattle.Opponent.GetActivePokemon(),
		Move:     move,
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) performSwitch(partyIndex int) {
	g.TrainerBattle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Switch,
		Selected: g.TrainerBattle.User.GetActivePokemon(),
		Target:   g.TrainerBattle.User.GetParty()[partyIndex],
	}
	g.SetActionBtns(Actions)
}
func (g *Gui) performUseItem(item *data.Item) {
	g.TrainerBattle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Bag,
		Selected: g.TrainerBattle.User.GetActivePokemon(),
		Target:   g.TrainerBattle.User.GetActivePokemon(),
		Item:     item,
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) renderActionBtns(gtx layout.Context) layout.Dimensions {
	switch g.TrainerBattle.ActionArea {
	case Actions:
		return g.renderBattleAction(gtx)
	case Attack:
		return g.renderBattleAttack(gtx)
	case EndBattle:
		return g.renderBattleEnd(gtx)
	default:
		return g.renderBattleAction(gtx)
	}
}
