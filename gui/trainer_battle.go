package gui

import (
	"log/slog"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"
)

type ActionAreaType int

const (
	Actions ActionAreaType = iota
	Attack
	EndBattle
)

type DialogActionAreaType int

const (
	MainBattle DialogActionAreaType = iota
	SwitchDialog
	BagDialog
)

const (
	AttackBtn       = "Attack"
	SwitchBtn       = "Switch"
	BagBtn          = "Bag"
	RunBtn          = "Run"
	Move1Btn        = "Move1"
	Move2Btn        = "Move2"
	Move3Btn        = "Move3"
	Move4Btn        = "Move4"
	EndBattleBtn    = "End Battle"
	SwitchDialogBtn = "Switch"
	CancelSwitchBtn = "Cancel"
)

type TrainerBattle struct {
	PokemonBattle battle.PokemonBattle
	User          battletrainer.BattleTrainer
	Opponent      battletrainer.BattleTrainer

	ActionChan chan data.BattleAction
	LogChan    chan string
	LogList    *widget.List
	LogContent []string

	ActionArea           ActionAreaType
	DialogActionArea     DialogActionAreaType
	ActionButtons        map[string]*widget.Clickable
	PokemonSwitchButtons []*widget.Clickable
	UseItemButtons       map[data.ItemName]*widget.Clickable
}

func (g *Gui) NewTrainerBattle(user *data.User, opponent *data.Trainer) TrainerBattle {
	battleAction := make(chan data.BattleAction, 1)
	logChan := make(chan string, 5)

	battleUser := battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{
		UserManager: g.opts.UserManager,
		PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{
			Logger: logger.NewChannelLogger(logChan),
		}),
	}, user, battleAction, logChan)

	battleTrainer := battletrainer.NewBattleAI(battletrainer.BattleTrainerOpts{
		UserManager:    nil,
		PokemonService: g.opts.PokemonService,
	}, opponent)

	pokemonBattle := battle.TrainerBattleManager(battle.TrainerBattleOpts{UserManager: g.opts.UserManager}, battleUser, battleTrainer)

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
	}
	pokemonSwitchBtns := make([]*widget.Clickable, len(battleUser.GetParty()))
	for i := range pokemonSwitchBtns {
		pokemonSwitchBtns[i] = new(widget.Clickable)
	}

	useItemBtns := make(map[data.ItemName]*widget.Clickable)
	for itemName := range user.Bag {
		useItemBtns[itemName] = new(widget.Clickable)
	}

	return TrainerBattle{
		PokemonBattle: pokemonBattle,
		Opponent:      battleTrainer,
		User:          battleUser,

		ActionChan: battleAction,
		LogChan:    logChan,
		LogList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		LogContent: make([]string, 0),

		ActionArea:           Actions,
		DialogActionArea:     MainBattle,
		ActionButtons:        buttons,
		PokemonSwitchButtons: pokemonSwitchBtns,
		UseItemButtons:       useItemBtns,
	}
}

func (g *Gui) LoadBattle(gtx layout.Context) layout.Dimensions {
	slog.Info("Starting Battle Screen")
	g.SetCurrentScreen(BattleScreen)

	go g.InitiateBattle()
	go g.LogBattle(gtx)

	return g.RenderBattleScreen(gtx)
}

func (g *Gui) InitiateBattle() {
	err := g.TrainerBattle.PokemonBattle.Introduction()
	if err != nil {
		slog.Error("Error found within battle", "error", err)
		panic(err)
	}

	g.SetActionBtns(EndBattle)
}

func (g *Gui) LogBattle(gtx layout.Context) {
	timeOutFlag := false
	for {
		select {
		case data := <-g.TrainerBattle.LogChan:
			slog.Info(data)
			g.TrainerBattle.LogContent = append(g.TrainerBattle.LogContent, data)
		case <-time.After(60 * time.Second):
			slog.Info("Timeout waiting for channel")
			timeOutFlag = true
		}

		if timeOutFlag {
			break
		}
	}
}

func (g *Gui) SetActionBtns(actionArea ActionAreaType) {
	g.TrainerBattle.ActionArea = actionArea
}
