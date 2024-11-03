package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sanchitdeora/PokeSim/data"
)

func (opts *GuiOpts) GetTrainerSideBar(user *data.User) fyne.CanvasObject {
	trainerSideBar := container.NewBorder(nil, nil, nil, nil, container.NewVBox(
		opts.GetTrainerInfo(user),
		widget.NewSeparator(),
		opts.GetBagInfo(user),
		widget.NewSeparator(),
		opts.GetPokemonInfo(user),
	))

	return addBorder(trainerSideBar)
}

func (opts *GuiOpts) GetTrainerInfo(user *data.User) fyne.CanvasObject {
	winPercentage := 0.0

	if user.Stats.Battles > 0 {
		winPercentage = (float64(user.Stats.Wins) / float64(user.Stats.Battles))
	}

	return container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Name:\t\t\t%s", user.Name)),
		widget.NewLabel(fmt.Sprintf("Money:\t\t\t$%v", user.Money)),
		widget.NewLabel(fmt.Sprintf("Catches:\t\t\t%v", user.Stats.Catches)),
		widget.NewLabel(fmt.Sprintf("Battles:\t\t\t%v", user.Stats.Battles)),
		widget.NewLabel(fmt.Sprintf("Wins:\t\t\t\t%v", user.Stats.Wins)),
		widget.NewLabel(fmt.Sprintf("Loss/Flees:\t\t\t%v", user.Stats.Losses)),
		widget.NewLabel(fmt.Sprintf("Win Percentage:\t\t%0.2f%%", winPercentage)),
		widget.NewLabel(fmt.Sprintf("PokéDEX:\t\t\t%v of 1000", user.Stats.PokeDEX)),
	)
}

func (opts *GuiOpts) GetBagInfo(user *data.User) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Pokeballs:\t\t\t%v", user.Bag[data.PokeBall].Count)),
		widget.NewLabel(fmt.Sprintf("Superballs:\t\t\t%v", user.Bag[data.SuperBall].Count)),
		widget.NewLabel(fmt.Sprintf("Ultraballs:\t\t\t%v", user.Bag[data.UltraBall].Count)),
		widget.NewLabel(fmt.Sprintf("Masterballs:\t\t%v", user.Bag[data.MasterBall].Count)),
		widget.NewLabel(fmt.Sprintf("Potion:\t\t\t%v", user.Bag[data.Potion].Count)),
		widget.NewLabel(fmt.Sprintf("Super Potion:\t\t%v", user.Bag[data.SuperPotion].Count)),
		widget.NewLabel(fmt.Sprintf("Hyper Potion:\t\t%v", user.Bag[data.HyperPotion].Count)),
	)
}

func (opts *GuiOpts) GetPokemonInfo(user *data.User) fyne.CanvasObject {
	return widget.NewLabel("Add Pokemon Here")
}
