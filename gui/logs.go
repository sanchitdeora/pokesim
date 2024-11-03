package gui

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (opts *GuiOpts) NewLogContainer() *fyne.Container {
	return container.NewVBox(widget.NewLabel("Battle Log:"))
}

func (opts *GuiOpts) LogListener() {
	slog.Info("Starting battle log listener", "from channel", opts.BattleLogChan)

	// Use a loop to keep listening for log messages
	for log := range opts.BattleLogChan {
		slog.Info("Battle log received", "log", log)
		opts.AppendLogContent(widget.NewLabel(log))
	}
	slog.Info("Battle log channel closed") // Optional: log closure of channel
}

// UpdateActionContent updates the actionContent container with new content
func (opts *GuiOpts) AppendLogContent(newContent fyne.CanvasObject) {
	opts.LogContainer.Add(newContent)
	opts.LogContainer.Refresh()
}

func (opts *GuiOpts) ClearLogContent() {
	opts.LogContainer.Objects = []fyne.CanvasObject{}
	opts.LogContainer.Refresh()
}
