package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)



func (opts *GuiOpts) DefaultActionTab() *fyne.Container {
	actionWidget := widget.NewLabel("Main Action Area")

	return container.NewBorder(
		actionWidget, nil, nil, nil,
	)
}

// UpdateActionContent updates the actionContent container with new content
func (opts *GuiOpts) UpdateActionContent(newContent fyne.CanvasObject) {
	opts.ActionContainer.Objects = []fyne.CanvasObject{newContent}
	opts.ActionContainer.Refresh()
}
