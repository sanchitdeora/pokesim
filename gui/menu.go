package gui

import (
	"fyne.io/fyne/v2"
	battle "github.com/sanchitdeora/PokeSim/battle_v1"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
)

func (opts *GuiOpts) MainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		&fyne.Menu{
			Label: "File",
			Items: []*fyne.MenuItem{
				{
					Label: "Save Game",
				},
				{
					Label: "Load Game",
				},
			},
		},
		&fyne.Menu{
			Label: "Battle",
			Items: []*fyne.MenuItem{
				{
					Label: "Test Trainer Battle",
					Action: func() {
						battleChan := make(chan *data.BattleInput, 1)
						battleLogChan := make(chan string, 1)
						opts.BattleLogChan = battleLogChan

						// create a new Battle Arena
						battleArena := NewBattleArena(opts,
							battleChan,
							battle.NewTrainerBattle(&battle.TrainerBattleOpts{
								UserService:     opts.UserService,
								PokemonService:  opts.PokemonService,
								BattleInputChan: battleChan,
								BattleLogChan:   battleLogChan,
							},
								trainermanagement.NewTrainerManager(trainermanagement.TrainerOpts{SavedTrainerPath: "C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\test_trainer.json"}).GetTrainer(),
							),
						)

						opts.UpdateActionContent(battleArena.LoadBattleScreen())
					},
				},
				{
					Label: "Wild Pokemon Battle",
					// Action: func () {opts.UpdateActionContent(opts.LoadBattleScreen()))},
				},
			},
		},
		&fyne.Menu{
			Label: "Trainer",
			Items: []*fyne.MenuItem{
				{
					Label: "PokeDex",
				},
				{
					Label: "Trainer Stats",
				},
				{
					Label: "Pokemon Box",
				},
			},
		},
		&fyne.Menu{
			Label: "PokeMart",
			Items: []*fyne.MenuItem{
				{
					Label: "Buy",
				},
				{
					Label: "Sell",
				},
			},
		},
	)
}
