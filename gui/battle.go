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
	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func (g *Gui) RenderBattleScreen(gtx layout.Context) layout.Dimensions {
	mainScreen := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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

	switch g.Battle.DialogActionArea {
	case SwitchDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderSwitchDialog(gtx)
					}),
				)
			}),
		)
	case BagDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderBagDialog(gtx)
					}),
				)
			}),
		)
	case EvolveDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderEvolveDialog(gtx)
					}),
				)
			}),
		)
	case LearnMoveDialog:
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return mainScreen
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return g.renderLearnMoveDialog(gtx)
					}),
				)
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
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(g.Battle.Opponent.GetActivePokemon().Name))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Lv. %v", g.Battle.Opponent.GetActivePokemon().Level))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// HP Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderHPBar(gtx, g.GetRemainingHPFraction(g.Battle.Opponent.GetActivePokemon())) // Example: 50% HP remaining
				})
			}),
		)
	})
}

// Opponent Image
func (g *Gui) renderOpponentImage(gtx layout.Context) layout.Dimensions {
	var img widget.Image

	if battle.IsPokemonInParty(g.Battle.User.GetParty(), g.Battle.Opponent.GetActivePokemon()) {
		img = loadImage("./assets/utils/pokeball.png")
		g.SetActionBtns(EndBattle)
	} else {
		img = loadImage(g.Battle.Opponent.GetActivePokemon().SpritesURL.FrontPath) // Replace with your image loading logic
	}
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
						label := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(g.Battle.User.GetActivePokemon().Name))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.Theme, fmt.Sprintf("Lv. %v", g.Battle.User.GetActivePokemon().Level))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// HP Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderHPBar(gtx, g.GetRemainingHPFraction(g.Battle.User.GetActivePokemon())) // Example: 50% HP remaining
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
						label := material.Body2(g.Theme, fmt.Sprintf("%d/%d HP", g.Battle.User.GetActivePokemon().BattleHP, g.Battle.User.GetActivePokemon().Stats.HP.Value))
						label.Alignment = text.Middle
						label.Font.Weight = font.Bold
						return label.Layout(gtx)
					}),
				)
			}),
			// Experience Bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(75), Right: unit.Dp(75)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.renderEXPBar(gtx, g.GetRemainingEXPFraction(g.Battle.User.GetActivePokemon().Pokemon))
				})
			}),
		)
	})
}

// User Image
func (g *Gui) renderUserImage(gtx layout.Context) layout.Dimensions {
	img := loadImage(g.Battle.User.GetActivePokemon().SpritesURL.BackPath) // Replace with your image loading logic
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
		Axis: layout.Vertical,
	}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
		var children []layout.FlexChild

		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.Battle.ActionButtons[AttackBtn].Clicked(gtx) {
							g.SetActionBtns(Attack)
						}
						btnStyle := material.Button(g.Theme, g.Battle.ActionButtons[AttackBtn], "Attack")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.Battle.ActionButtons[SwitchBtn].Clicked(gtx) {
							g.Battle.DialogActionArea = SwitchDialog
						}

						btnStyle := material.Button(g.Theme, g.Battle.ActionButtons[SwitchBtn], "Switch")
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
				Axis: layout.Horizontal,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.Battle.ActionButtons[BagBtn].Clicked(gtx) {
							slog.Info("Bag Button Clicked")
							g.Battle.DialogActionArea = BagDialog
						}

						btnStyle := material.Button(g.Theme, g.Battle.ActionButtons[BagBtn], "Bag")
						btnStyle.Background = SecondaryBackgroundColor
						btnStyle.Color = TextColor

						return btnStyle.Layout(gtx)
					})
				}))
				children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if g.Battle.ActionButtons[RunBtn].Clicked(gtx) {
							g.performRun()
							g.SetActionBtns(Actions)
						}

						btnStyle := material.Button(g.Theme, g.Battle.ActionButtons[RunBtn], "Run")
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
	activePokemon := g.Battle.User.GetActivePokemon()

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
		var children []layout.FlexChild

		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children,
					g.renderMoveBtn(gtx, activePokemon.Moveset.Move1, g.Battle.ActionButtons[Move1Btn]),
					g.renderMoveBtn(gtx, activePokemon.Moveset.Move2, g.Battle.ActionButtons[Move2Btn]),
				)
				return children
			}(gtx)...,
			)
		}))
		children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
				var children []layout.FlexChild

				children = append(children,
					g.renderMoveBtn(gtx, activePokemon.Moveset.Move3, g.Battle.ActionButtons[Move3Btn]),
					g.renderMoveBtn(gtx, activePokemon.Moveset.Move4, g.Battle.ActionButtons[Move4Btn]),
				)
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
								}).Layout(gtx, len(g.Battle.User.GetParty()), func(gtx layout.Context, i int) layout.Dimensions {
									pokemon := g.Battle.User.GetParty()[i]

									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										// Pokémon Image
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											img := loadImage(pokemon.SpritesURL.FrontPath)
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
												btnStyle := material.Button(g.Theme, g.Battle.PokemonSwitchButtons[i], "Switch")
												if g.Battle.PokemonSwitchButtons[i].Clicked(gtx) {
													slog.Info("Switching to Pokémon", "pokemon", pokemon.Name)
													if g.Battle.User.GetParty()[i].BattleHP > 0 {
														g.performSwitch(i)
														g.Battle.DialogActionArea = MainBattle
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
								btn := material.Button(g.Theme, g.Battle.ActionButtons[CancelSwitchBtn], "Cancel")
								btn.Background = RedBtnColor
								if g.Battle.ActionButtons[CancelSwitchBtn].Clicked(gtx) {
									g.Battle.DialogActionArea = MainBattle
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

func (g *Gui) canUseItemInBattle(item data.Item) bool {
	return g.Battle.CatchPokemonEnabled || item.Category != data.PokeBalls
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
									item, exists := g.Battle.User.GetTrainer().Bag[itemName]
									if !exists || !g.canUseItemInBattle(item) || item.Count < 1 {
										return layout.Dimensions{}
									}

									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										// Item Image
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
											label := material.Body1(g.Theme, fmt.Sprintf("x%v", item.Count))
											label.Alignment = text.Middle
											// label.Font.Weight = font.Bold
											return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return label.Layout(gtx)
											})

										}),
										layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
											return layout.Inset(layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												btnStyle := material.Button(g.Theme, g.Battle.UseItemButtons[itemName], "Use")
												if g.Battle.UseItemButtons[itemName].Clicked(gtx) {
													if (g.Battle.User.GetTrainer().Bag[itemName].Category == data.MedicalItems && (g.Battle.User.GetActivePokemon().BattleHP < g.Battle.User.GetActivePokemon().Pokemon.Stats.HP.Value)) ||
														g.Battle.User.GetTrainer().Bag[itemName].Category == data.PokeBalls && (g.Battle.Opponent.GetActivePokemon().BattleHP > 0) {

														g.performUseItem(&item)
														g.Battle.DialogActionArea = MainBattle
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
								btn := material.Button(g.Theme, g.Battle.ActionButtons[CancelSwitchBtn], "Cancel")
								btn.Background = RedBtnColor
								if g.Battle.ActionButtons[CancelSwitchBtn].Clicked(gtx) {
									g.Battle.DialogActionArea = MainBattle
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
		if g.Battle.ActionButtons[EndBattleBtn].Clicked(gtx) {
			g.SetCurrentScreen(HomeScreen)
		}

		btnStyle := material.Button(g.Theme, g.Battle.ActionButtons[EndBattleBtn], "End Battle")
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
				return material.List(g.Theme, g.Battle.LogList).Layout(gtx, len(g.Battle.LogContent), func(gtx layout.Context, index int) layout.Dimensions {
					gtx.Constraints.Max.Y = availableHeight
					return material.Body2(g.Theme, g.Battle.LogContent[index]).Layout(gtx)
				})
			})
		}),
	)
}

func (g *Gui) renderMoveBtn(gtx layout.Context, move *data.Moves, moveBtn *widget.Clickable) layout.FlexChild {
	isDisabled := false
	if move == nil || move.Name == "" {
		isDisabled = true
		move = &data.Moves{Name: ""}
	}

	return layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if !isDisabled && moveBtn.Clicked(gtx) {
				g.performAttack(move)
			}

			btnStyle := material.Button(g.Theme, moveBtn, utils.ToCapitalizeFirstLetterOfEachWord(move.Name))
			if isDisabled {
				btnStyle.Background = ButtonDisabledColor
			} else {
				btnStyle.Background = SecondaryBackgroundColor
			}
			btnStyle.Color = TextColor

			return btnStyle.Layout(gtx)
		})
	})
}

func (g *Gui) renderLearnNewMoveSelection(gtx layout.Context, index int, move data.Moves) layout.FlexChild {
	return layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)

			slog.Info("renderLearnNewMoveSelection", "index", index, "selectedIndex", g.Battle.LearnNewMove.SelectedIndex)

			isSelected := g.Battle.LearnNewMove.SelectedIndex == index
			bgColor := PrimaryBackgroundColor // Default background
			if isSelected {
				bgColor = SelectedColor // Highlighted background
			}

			drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
			fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, bgColor, 8)

			return layout.Stack{
				Alignment: layout.Center,
			}.Layout(gtx,
				// Move Name
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Spacing:   layout.SpaceBetween,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(move.Name))
							name.Alignment = text.Middle
							return name.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							name := material.Body2(g.Theme, fmt.Sprintf("Type: %s", utils.ToCapitalizeFirstLetterOfEachWord(string(move.Type))))
							name.Alignment = text.Middle
							return name.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Spacing:   layout.SpaceBetween,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									name := material.Body2(g.Theme, fmt.Sprintf("Power: %v", move.Power))
									name.Alignment = text.Middle
									return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return name.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									name := material.Body2(g.Theme, fmt.Sprintf("Accuracy: %v", move.Accuracy))
									name.Alignment = text.Middle
									return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return name.Layout(gtx)
									})
								}),
							)
						}),
					)
					// name := material.Body1(g.Theme, utils.ToCapitalizeFirstLetterOfEachWord(move.Name))
					// name.Alignment = text.Middle
					// return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// 	return name.Layout(gtx)
					// })
				}),
				// Clickable Overlay for selection
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					if g.Battle.LearnNewMove.ForgetMovesBtns[index].Clicked(gtx) {
						slog.Info("clickable overlay", "index", index, "btns len", len(g.Battle.LearnNewMove.ForgetMovesBtns), "btns", g.Battle.LearnNewMove.ForgetMovesBtns)
						g.Battle.LearnNewMove.SelectedIndex = index
					}
					return g.Battle.LearnNewMove.ForgetMovesBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: cardDims}
					})
				}),
			)
		})
	})
}

func (g *Gui) performAttack(move *data.Moves) {
	g.Battle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: g.Battle.User.GetActivePokemon(),
		Target:   g.Battle.Opponent.GetActivePokemon(),
		Move:     move,
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) performSwitch(partyIndex int) {
	g.Battle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Switch,
		Selected: g.Battle.User.GetActivePokemon(),
		Target:   g.Battle.User.GetParty()[partyIndex],
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) performUseItem(item *data.Item) {
	var target *data.BattlePokemon
	if item.Category == data.PokeBalls {
		target = g.Battle.Opponent.GetActivePokemon()
	} else {
		target = g.Battle.User.GetActivePokemon()
	}

	g.Battle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Bag,
		Selected: g.Battle.User.GetActivePokemon(),
		Target:   target,
		Item:     item,
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) performRun() {
	g.Battle.ActionChan <- data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Run,
		Selected: g.Battle.User.GetActivePokemon(),
		Target:   g.Battle.Opponent.GetActivePokemon(),
	}
	g.SetActionBtns(Actions)
}

func (g *Gui) renderActionBtns(gtx layout.Context) layout.Dimensions {
	switch g.Battle.ActionArea {
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

func (g *Gui) renderEvolveDialog(gtx layout.Context) layout.Dimensions {
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
								label := material.H6(g.Theme, "Evolve Pokemon?")
								label.Alignment = text.Middle
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							// Pokémon list
							layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
								evolveBody := g.Battle.LevelUpBody.(data.EventEvolveBody)

								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,

									// Base Pokémon Image
									layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
										img := loadImage(evolveBody.CurrentBasePokemon.SpritesURL.FrontPath)
										img.Fit = widget.Contain
										img.Position = layout.Center
										return img.Layout(gtx)
									}),
									// Evo Arrow
									layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
										img := loadImage("assets/trainer/img/new_game_img.png")
										img.Fit = widget.Contain
										img.Position = layout.Center
										return img.Layout(gtx)
									}),
									// Pokémon Image
									layout.Flexed(1.0, func(gtx layout.Context) layout.Dimensions {
										var evolvedPokemonFlexChildren []layout.FlexChild

										for i, basePokemon := range evolveBody.EvolvedBasePokemons {

											isSelected := g.Battle.EvolutionProps.SelectedEvolveBasePokemon == &evolveBody.EvolvedBasePokemons[i]
											bgColor := SecondaryBackgroundColor // Default background
											if isSelected {
												bgColor = SelectedColor // Highlighted background
											}

											if g.Battle.EvolutionProps.EvolvedPokemonSelectedBtns[i] != nil && g.Battle.EvolutionProps.EvolvedPokemonSelectedBtns[i].Clicked(gtx) {
												g.Battle.EvolutionProps.SelectedEvolveBasePokemon = &evolveBody.EvolvedBasePokemons[i]
												slog.Info("evo pokemon Selected Clicked", "index", i, "selectedPokemon", g.Battle.EvolutionProps.SelectedEvolveBasePokemon.Name)
											}

											evolvedPokemonFlexChildren = append(evolvedPokemonFlexChildren,
												layout.Flexed(1.0/float32(len(evolveBody.EvolvedBasePokemons)), func(gtx layout.Context) layout.Dimensions {
													cardDims := gtx.Constraints.Max

													drawRoundedBorder(gtx, cardDims, CardBorderColor, unit.Dp(1), unit.Dp(8))
													fillRoundedShape(gtx, image.Rectangle{Max: cardDims}, bgColor, 8)

													return layout.Stack{Alignment: layout.Center}.Layout(gtx,
														layout.Stacked(func(gtx layout.Context) layout.Dimensions {
															img := loadImage(basePokemon.SpritesURL.FrontPath)
															img.Fit = widget.Contain
															img.Position = layout.Center
															return img.Layout(gtx)
														}),
														layout.Stacked(func(gtx layout.Context) layout.Dimensions {
															if g.Battle.EvolutionProps.EvolvedPokemonSelectedBtns[i] == nil {
																return layout.Dimensions{Size: cardDims}
															}
															return g.Battle.EvolutionProps.EvolvedPokemonSelectedBtns[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
																return layout.Dimensions{Size: cardDims}
															})
														}),
													)
												}),
											)
										}

										return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx, evolvedPokemonFlexChildren...)
									}),
								)
							}),

							// Buttons
							layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {

								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
									layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
										isBtnDisabled := g.Battle.EvolutionProps.SelectedEvolveBasePokemon == nil

										btnBgColor := SecondaryBackgroundColor // Grey background for disabled
										if isBtnDisabled {
											btnBgColor = ButtonDisabledColor // Active background (blue)
										}

										if g.Battle.EvolutionProps.EvolveBtn.Clicked(gtx) && !isBtnDisabled {
											g.SendLevelUpResponse(gtx, data.LevelUpEvent{
												EventType: data.LevelUpEventEvolve,
												Body: data.ResponseEvolveBody{
													AcceptEvolution:    true,
													EvolvedBasePokemon: *g.Battle.EvolutionProps.SelectedEvolveBasePokemon,
												}},
											)
											g.Battle.DialogActionArea = MainBattle
											g.SetActionBtns(EndBattle)
											g.Battle.EvolutionProps.SelectedEvolveBasePokemon = nil
										}
										if g.Battle.EvolutionProps.EvolveBtn.Hovered() {
											btnBgColor = ButtonHoveredColor
										}

										btn := material.Button(g.Theme, g.Battle.EvolutionProps.EvolveBtn, "Evolve!")
										btn.Color = TextColor
										btn.Background = btnBgColor
										return layout.Center.Layout(gtx, btn.Layout)
									}),
									layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(g.Theme, g.Battle.EvolutionProps.CancelEvolveBtn, "Cancel")
										btn.Background = RedBtnColor

										if g.Battle.EvolutionProps.CancelEvolveBtn.Clicked(gtx) {
											g.SendLevelUpResponse(gtx, data.LevelUpEvent{
												EventType: data.LevelUpEventEvolve,
												Body: data.ResponseEvolveBody{
													AcceptEvolution: false,
												}},
											)
											g.Battle.DialogActionArea = MainBattle
											g.SetActionBtns(EndBattle)
										}
										return layout.Center.Layout(gtx, btn.Layout)
									}),
								)
							}),
						)
					}),
				)
			})
		}),
	)
}

func (g *Gui) renderLearnMoveDialog(gtx layout.Context) layout.Dimensions {
	learnMoveBody := g.Battle.LevelUpBody.(data.EventLearnMoveBody)

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
								label := material.H6(g.Theme, "Learn Move?")
								label.Alignment = text.Middle
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							// Move list
							layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,

									layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
										label := material.H6(g.Theme, fmt.Sprintf("%s wants a learn a new move. Which move should be forgotten?", utils.ToCapitalizeFirstLetterOfEachWord(learnMoveBody.Pokemon.Name)))
										label.Alignment = text.Middle
										label.Font.Weight = font.Bold
										return label.Layout(gtx)
									}),
									// Pokémon Image
									layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
											var children []layout.FlexChild

											children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{
													Axis: layout.Horizontal,
												}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
													var children []layout.FlexChild

													children = append(children,
														g.renderLearnNewMoveSelection(gtx, 1, *learnMoveBody.Pokemon.Moveset.Move1),
														g.renderLearnNewMoveSelection(gtx, 2, *learnMoveBody.Pokemon.Moveset.Move2),
													)
													return children
												}(gtx)...,
												)
											}))
											children = append(children, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{
													Axis: layout.Horizontal,
												}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
													var children []layout.FlexChild

													children = append(children,
														g.renderLearnNewMoveSelection(gtx, 3, *learnMoveBody.Pokemon.Moveset.Move3),
														g.renderLearnNewMoveSelection(gtx, 4, *learnMoveBody.Pokemon.Moveset.Move4),
													)
													return children
												}(gtx)...,
												)
											}))
											return children
										}(gtx)...,
										)
									}),
									// New Move
									layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis:    layout.Horizontal,
											Spacing: layout.SpaceBetween,
										}.Layout(gtx, func(gtx layout.Context) []layout.FlexChild {
											var children []layout.FlexChild

											children = append(children,
												g.renderLearnNewMoveSelection(gtx, 0, learnMoveBody.NewMove),
											)
											return children
										}(gtx)...,
										)
									}),
								)
							}),

							// Buttons
							layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
								// initialized early to get accurate size
								return layout.Flex{
									Axis:    layout.Horizontal,
									Spacing: layout.SpaceBetween,
								}.Layout(gtx,
									layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {

										btnBgColor := PrimaryBackgroundColor
										if g.Battle.LearnNewMove.SelectedIndex == -1 {
											btnBgColor = ButtonDisabledColor
										}

										if g.Battle.ActionButtons[LearnMoveBtn].Clicked(gtx) {
											g.SendLearnNewMoveResponse(gtx, g.Battle.LearnNewMove.SelectedIndex, learnMoveBody)
											g.Battle.LearnNewMove.SelectedIndex = -1
											g.Battle.DialogActionArea = MainBattle
										}
										if g.Battle.ActionButtons[LearnMoveBtn].Hovered() {
											btnBgColor = ButtonHoveredColor
										}

										return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											btnSize := image.Point{X: 120, Y: 60}
											radius := 5

											// // define rounded border and fill
											// drawRoundedBorder(gtx, btnSize, PrimaryBackgroundColor, unit.Dp(1), unit.Dp(radius))
											// fillRoundedShape(gtx, image.Rectangle{Max: btnSize}, btnBgColor, radius)

											return layout.Stack{
												Alignment: layout.Center,
											}.Layout(gtx,
												layout.Stacked(func(gtx layout.Context) layout.Dimensions {
													// Draw the border around the text
													drawRoundedBorder(gtx, btnSize, PrimaryBackgroundColor, unit.Dp(1), unit.Dp(radius))
													fillRoundedShape(gtx, image.Rectangle{Max: btnSize}, btnBgColor, radius)
													return layout.Dimensions{Size: btnSize}
												}),
												layout.Stacked(func(gtx layout.Context) layout.Dimensions {
													textStyle := material.Body2(g.Theme, "Forget Move!")
													textStyle.Alignment = text.Middle
													textStyle.Color = TextColor

													return textStyle.Layout(gtx)
												}),
												layout.Expanded(func(gtx layout.Context) layout.Dimensions {
													if g.Battle.LearnNewMove.SelectedIndex == -1 {
														return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
															return layout.Dimensions{Size: btnSize}
														})
													}

													return g.Battle.ActionButtons[LearnMoveBtn].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
														return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
															return layout.Dimensions{Size: btnSize}
														})
													})
												}),
											)
										})
									}),
									layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(g.Theme, g.Battle.ActionButtons[CancelSwitchBtn], "Cancel")
										btn.Background = RedBtnColor
										if g.Battle.ActionButtons[CancelSwitchBtn].Clicked(gtx) {
											g.SendLearnNewMoveResponse(gtx, 0, learnMoveBody)
											g.Battle.LearnNewMove.SelectedIndex = -1
											g.Battle.DialogActionArea = MainBattle
										}
										return layout.Center.Layout(gtx, btn.Layout)
									}),
								)
							}),
						)
					}),
				)
			})
		}),
	)
}

func (g *Gui) SendLearnNewMoveResponse(gtx layout.Context, selectedIndex int, learnMove data.EventLearnMoveBody) {
	if selectedIndex == -1 {
		slog.Warn("Selected Index is -1")
		return
	}

	var replacedMove data.Moves
	updatedMoveset := learnMove.Pokemon.Moveset

	if selectedIndex == 0 {
		replacedMove = learnMove.NewMove
	} else if selectedIndex == 1 {
		replacedMove = *learnMove.Pokemon.Moveset.Move1
		updatedMoveset.Move1 = &learnMove.NewMove
	} else if selectedIndex == 2 {
		replacedMove = *learnMove.Pokemon.Moveset.Move2
		updatedMoveset.Move2 = &learnMove.NewMove
	} else if selectedIndex == 3 {
		replacedMove = *learnMove.Pokemon.Moveset.Move3
		updatedMoveset.Move3 = &learnMove.NewMove
	} else if selectedIndex == 4 {
		replacedMove = *learnMove.Pokemon.Moveset.Move4
		updatedMoveset.Move4 = &learnMove.NewMove
	}

	g.SendLevelUpResponse(gtx, data.LevelUpEvent{
		EventType: data.LevelUpEventLearnMove,
		Body: data.ResponseLearnMoveBody{
			UpdatedMoveset: updatedMoveset,
			NewMove:        learnMove.NewMove,
			ReplacedMove:   replacedMove,
		}},
	)
}

type EvolutionProps struct {
	SelectedEvolveBasePokemon  *data.BasePokemon
	EvolvedPokemonSelectedBtns []*widget.Clickable
	EvolveBtn                  *widget.Clickable
	CancelEvolveBtn            *widget.Clickable
}
