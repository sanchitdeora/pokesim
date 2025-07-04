package gui

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func DefaultTrainerProps() TrainerProps {
	return TrainerProps{
		TrainerList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		Trainers: SetupTrainersList(),
	}
}

func DefaultGymTrainerProps() TrainerProps {
	return TrainerProps{
		TrainerList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		Trainers: SetupGymTrainersList(),
	}
}

type TrainerProps struct {
	TrainerList *widget.List
	Trainers    []TrainerUI
}

func (g *Gui) RenderTrainerScreen(gtx layout.Context) layout.Dimensions {
	// Use a vertical layout for screen elements
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Trainers")
			title.Font.Weight = font.Bold
			return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					trainers := g.TrainerProps.Trainers
					return g.renderTrainerGallery(gtx, trainers)
				}))
		}),
	)
}

func (g *Gui) RenderGymTrainerScreen(gtx layout.Context) layout.Dimensions {
	// Use a vertical layout for screen elements
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Gym Trainers")
			title.Font.Weight = font.Bold
			return layout.Inset(layout.Inset{Bottom: unit.Dp(25)}).Layout(gtx, func(gtx layout.Context) layout.Dimensions { return title.Layout(gtx) })
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gymTrainers := g.GymTrainerProps.Trainers
					return g.renderTrainerGallery(gtx, gymTrainers)
				}))
		}),
	)
}

func (g *Gui) renderTrainerGallery(gtx layout.Context, trainers []TrainerUI) layout.Dimensions {
	// Wrap the trainer list in a flex layout for better control

	return layout.Flex{}.Layout(gtx,
		// Expanding the list container to fill available height
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(g.Theme, g.TrainerProps.TrainerList).Layout(gtx, len(trainers), func(gtx layout.Context, index int) layout.Dimensions {
				return g.renderTrainerRow(gtx, trainers, index)
			})
		}),
	)
}

func (g *Gui) renderTrainerRow(gtx layout.Context, trainers []TrainerUI, index int) layout.Dimensions {
	// Create a row with up to 4 trainers
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		func(gtx layout.Context) []layout.FlexChild {
			var children []layout.FlexChild
			for i := 0; i < 4; i++ {
				idx := index*4 + i
				if idx >= len(trainers) {
					// Add placeholder for missing trainers
					children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
						// return g.renderPlaceholderCard(gtx)
					}))
					continue
				}
				trainer := trainers[idx]
				children = append(children, layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
					return g.renderTrainerCard(gtx, trainer)
				}))
			}
			return children
		}(gtx)...,
	)
}

func (g *Gui) renderTrainerCard(gtx layout.Context, trainer TrainerUI) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		cardDims := image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(250)))

		// Paint the rounded border
		cardColor := TrainerLockedColor
		if trainer.Unlocked {
			cardColor = CardBorderColor
		}
		drawRoundedBorder(gtx, cardDims, cardColor, unit.Dp(1), unit.Dp(8))

		return layout.Stack{
			Alignment: layout.Center,
		}.Layout(gtx,
			// Trainer Image
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				img := trainer.Image
				newCtx := gtx
				newCtx.Constraints.Max = cardDims
				newCtx.Constraints.Max.Y -= 50
				if !trainer.Unlocked {
					// Grey out the image for locked trainers

					return g.greyedOutImage(newCtx, img)
				}
				img.Fit = widget.Contain
				return img.Layout(newCtx)
			}),
			// Trainer Name
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				name := material.Body1(g.Theme, trainer.Trainer.Name)
				name.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(cardDims.Y - int(name.TextSize) - 20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return name.Layout(gtx)
				})
			}),
			// Clickable Overlay for unlocked trainers
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				if trainer.Unlocked {
					if trainer.TrainerBtn.Clicked(gtx) {
						g.Battle = g.NewTrainerBattle(g.opts.UserManager.GetUser(), trainer.Trainer)
						return g.LoadBattle(gtx)
					}

					return trainer.TrainerBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: cardDims}
					})
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func (g *Gui) greyedOutImage(gtx layout.Context, img widget.Image) layout.Dimensions {
	// Render the image first
	imageDims := img.Layout(gtx)

	// Apply a translucent grey overlay
	greyColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128} // Semi-transparent grey
	fillRoundedShape(gtx, image.Rectangle{Max: imageDims.Size}, greyColor, 128)

	return imageDims
}

var TrainerLockedColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Light grey for locked trainers
var TrainerUnlockedColor = color.NRGBA{R: 34, G: 139, B: 34, A: 255} // Forest green for unlocked trainers
