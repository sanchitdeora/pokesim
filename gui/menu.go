package gui

import (
	"fyne.io/fyne/v2"
	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
)

func (opts *GuiOpts) MainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		&fyne.Menu{
			Label: "File",
			Items: []*fyne.MenuItem{
				{
					ChildMenu: &fyne.Menu{
						Label: "Save Game",
						Items: []*fyne.MenuItem{
							{},
						},
					},
					// Icon: theme.DocumentSaveIcon(),
				},
				{
					ChildMenu: &fyne.Menu{
						Label: "Load Game",
						Items: []*fyne.MenuItem{
							{},
						},
					},
					// Icon: theme.DocumentIcon(),
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

						opts.UpdateActionContent(
							LoadBattleScreen(opts,
								battleChan,
								battle.NewTrainerBattle(&battle.TrainerBattleOpts{
									UserService:     opts.UserService,
									PokemonService:  opts.PokemonService,
									BattleInputChan: battleChan,
									BattleLogChan:   battleLogChan,
								},
									trainermanagement.NewTrainer(trainermanagement.TrainerOpts{SavedTrainerPath: "C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\test_trainer.json"}).GetTrainer(),
								),
							),
						)
					},
					// Icon: theme.CancelIcon(),
				},
				{
					Label: "Wild Pokemon Battle",
					// Action: func () {opts.UpdateActionContent(opts.LoadBattleScreen()))},
				},
			},
		},
		&fyne.Menu{
			Label: "Shop",
			Items: []*fyne.MenuItem{
				{
					ChildMenu: &fyne.Menu{
						Label: "Buy",
						Items: []*fyne.MenuItem{
							{},
						},
					},
					// Icon: theme.DesktopIcon(),
				},
				{
					ChildMenu: &fyne.Menu{
						Label: "Sell",
						Items: []*fyne.MenuItem{
							{},
						},
					},
					// Icon: theme.CheckButtonIcon(),
				},
			},
		},
		&fyne.Menu{
			Label: "PokeDEX",
			Items: []*fyne.MenuItem{
				{
					ChildMenu: &fyne.Menu{
						Label: "Open PokeDex",
						Items: []*fyne.MenuItem{
							{},
						},
					},
					// Icon: theme.ComputerIcon(),
				},
			},
		},
	)
}
