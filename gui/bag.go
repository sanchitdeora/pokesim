package gui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (g *Gui) renderBagWindow(gtx layout.Context) layout.Dimensions {
	// Create a theme for styling
	th := material.NewTheme()

	// Use a vertical layout for screen elements
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Title
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(th, "Bag Screen")
			return title.Layout(gtx)
		}),
	)
}
