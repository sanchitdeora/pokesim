package gui

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (g *Gui) RenderBoxScreen(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,

		// Title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(g.Theme, "Box")
			title.Font.Weight = font.Bold

			return title.Layout(gtx)
		}),
	)
}
