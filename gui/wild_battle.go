package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"
)

func DefaultWildEnvironmentProps() WildEnvironmentProps {
	return WildEnvironmentProps{
		WildEnvironmentList:       &widget.List{
				List: layout.List{Axis: layout.Vertical},
			},
		Environments: getEnvironments(),
	}
}

type WildEnvironmentProps struct {
	WildEnvironmentList *widget.List
	Environments        []EnvironmentUI
}

func (g *Gui) NewWildBattle(user *data.User, wild *data.Pokemon) Battle {
	battleAction := make(chan data.BattleAction, 1)
	levelUpEventsChan := make(chan data.LevelUpEvent, 1)
	levelUpResponseChan := make(chan data.LevelUpEvent, 1)
	logChan := make(chan string, 5)

	battleUser := battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{
		UserManager: g.opts.UserManager,
		PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{
			Logger:           logger.NewChannelLogger(logChan),
			GameStateManager: g.opts.GameManager,
			LvlUpActions:     pokemon.NewLevelUpEvents(levelUpEventsChan, levelUpResponseChan),
		}),
	}, user, battleAction, logChan)

	battleTrainer := battletrainer.NewBattleWild(battletrainer.BattleTrainerOpts{
		UserManager:    nil,
		PokemonService: g.opts.PokemonService,
	}, wild)

	pokemonBattle := battle.WildBattleManager(battle.WildBattleOpts{UserManager: g.opts.UserManager}, battleUser, battleTrainer)

	buttons := map[string]*widget.Clickable{
		AttackBtn:       new(widget.Clickable),
		SwitchBtn:       new(widget.Clickable),
		BagBtn:          new(widget.Clickable),
		RunBtn:          new(widget.Clickable),
		Move1Btn:        new(widget.Clickable),
		Move2Btn:        new(widget.Clickable),
		Move3Btn:        new(widget.Clickable),
		Move4Btn:        new(widget.Clickable),
		EndBattleBtn:    new(widget.Clickable),
		CancelSwitchBtn: new(widget.Clickable),
		EvolveBtn:       new(widget.Clickable),
		LearnMoveBtn:    new(widget.Clickable),
	}

	pokemonSwitchBtns := make([]*widget.Clickable, len(battleUser.GetParty()))
	for i := range pokemonSwitchBtns {
		pokemonSwitchBtns[i] = new(widget.Clickable)
	}

	useItemBtns := make(map[data.ItemName]*widget.Clickable)
	for itemName := range user.Bag {
		useItemBtns[itemName] = new(widget.Clickable)
	}

	forgetMovesBtns := make([]*widget.Clickable, 5)
	for i := range forgetMovesBtns {
		forgetMovesBtns[i] = new(widget.Clickable)
	}

	return Battle{
		PokemonBattle: pokemonBattle,
		Opponent:      battleTrainer,
		User:          battleUser,

		CatchPokemonEnabled: true,
		RunBattleEnabled:    true,

		ActionChan: battleAction,

		LevelUpEventsChan:   levelUpEventsChan,
		LevelUpResponseChan: levelUpResponseChan,

		LogChan: logChan,
		LogList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		LogContent: make([]string, 0),

		ActionArea:           Actions,
		DialogActionArea:     MainBattle,
		ActionButtons:        buttons,
		PokemonSwitchButtons: pokemonSwitchBtns,
		UseItemButtons:       useItemBtns,

		LearnNewMove: LearnNewMoveUI{
			SelectedIndex:   -1,
			ForgetMovesBtns: forgetMovesBtns,
		},
	}
}
