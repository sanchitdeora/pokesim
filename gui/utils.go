package gui

import (
	"image"
	"image/color"
	"log/slog"
	"os"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

// HP Bar Renderer
func (g *Gui) renderHPBar(gtx layout.Context, percentage float32) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Render the HP label
			return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(g.Theme, "HP")
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Render the progress bar stretched across the remaining width
			barHeight := gtx.Dp(unit.Dp(8))              // Set bar height
			totalWidth := float32(gtx.Constraints.Max.X) // Total width available
			filledWidth := totalWidth * percentage       // Width based on progress

			paint.FillShape(gtx.Ops, color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, clip.Rect{
				Max: image.Pt(int(filledWidth), barHeight),
			}.Op())
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}, clip.Rect{
				Min: image.Pt(int(filledWidth), 0),
				Max: image.Pt(int(totalWidth), barHeight),
			}.Op())

			return layout.Dimensions{Size: image.Pt(int(totalWidth), barHeight)}
		}),
	)
}

// XP Bar Renderer
func (g *Gui) renderEXPBar(gtx layout.Context, percentage float32) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Render the HP label
			return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body2(g.Theme, "EXP").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Render the progress bar stretched across the remaining width
			barHeight := gtx.Dp(unit.Dp(8))              // Set bar height
			totalWidth := float32(gtx.Constraints.Max.X) // Total width available
			filledWidth := totalWidth * percentage       // Width based on progress

			paint.FillShape(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, clip.Rect{
				Max: image.Pt(int(filledWidth), barHeight),
			}.Op())
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}, clip.Rect{
				Min: image.Pt(int(filledWidth), 0),
				Max: image.Pt(int(totalWidth), barHeight),
			}.Op())

			return layout.Dimensions{Size: image.Pt(int(totalWidth), barHeight)}
		}),
	)
}

func (g *Gui) GetRemainingHPFraction(bPokemon *data.BattlePokemon) float32 {
	return float32(bPokemon.BattleHP) / float32(bPokemon.Stats.HP.Value)
}

func (g *Gui) GetRemainingEXPFraction(pokemon *data.Pokemon) float32 {
	totalExp := g.opts.PokemonService.GetExperienceRequiredForNextLevel(pokemon)

	// slog.Info("GetRemainingEXPFraction", "totalExp", totalExp, "bPokemon.ExperienceLeft", bPokemon.ExperienceLeft)

	return float32(totalExp-pokemon.ExperienceLeft) / float32(totalExp)
}

func drawRoundedBorder(gtx layout.Context, size image.Point, borderColor color.NRGBA, borderWidth unit.Dp, cornerRadius unit.Dp) {
	// Create a rounded rectangle
	r := gtx.Dp(cornerRadius)
	rrect := clip.RRect{
		Rect: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: size.X, Y: size.Y},
		},
		NE: r, NW: r, SE: r, SW: r,
	}

	// Draw the border by stroking the path
	strokePath := clip.Stroke{
		Path:  rrect.Path(gtx.Ops),
		Width: float32(gtx.Dp(borderWidth)),
	}.Op()

	paint.FillShape(gtx.Ops, borderColor, strokePath)
}

func fillRoundedShape(gtx layout.Context, size image.Rectangle, bgColor color.NRGBA, cornerRadius int) {
	paint.FillShape(gtx.Ops, bgColor, clip.RRect{Rect: size,
		NE: cornerRadius, NW: cornerRadius, SE: cornerRadius, SW: cornerRadius}.Op(gtx.Ops))
}

func loadImage(relativePath string) widget.Image {
	f, err := os.Open(utils.GetFullPath(relativePath)) // Adjust the path as needed
	if err != nil {
		slog.Info("Failed to load image", "error", err)
		return widget.Image{}
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		slog.Info("Failed to decode image", "error", err)
		return widget.Image{}
	}

	return widget.Image{
		Src: paint.NewImageOp(img),
		Fit: widget.Cover,
	}
}
