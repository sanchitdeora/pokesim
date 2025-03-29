package gui

import (
	"log/slog"

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
	EvolveDialog
	LearnMoveDialog
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
	NewMoveBtn      = "NewMove"
	EndBattleBtn    = "End Battle"
	SwitchDialogBtn = "Switch"
	CancelSwitchBtn = "Cancel"
	EvolveBtn       = "Evolve"
	LearnMoveBtn    = "LearnMove"
)

type Battle struct {
	PokemonBattle battle.PokemonBattle
	User          battletrainer.BattleTrainer
	Opponent      battletrainer.BattleTrainer

	CatchPokemonEnabled bool
	RunBattleEnabled   bool

	ActionChan chan data.BattleAction

	LevelUpEventsChan   chan data.LevelUpEvent
	LevelUpResponseChan chan data.LevelUpEvent
	LevelUpBody         interface{}

	LogChan    chan string
	LogList    *widget.List
	LogContent []string

	ActionArea           ActionAreaType
	DialogActionArea     DialogActionAreaType
	ActionButtons        map[string]*widget.Clickable
	PokemonSwitchButtons []*widget.Clickable
	UseItemButtons       map[data.ItemName]*widget.Clickable

	// Learn New Move Props
	LearnNewMove LearnNewMoveUI
}

func (g *Gui) NewTrainerBattle(user *data.User, opponent *data.Trainer) Battle {
	battleAction := make(chan data.BattleAction, 1)
	logChan := make(chan string, 5)
	levelUpEventsChan := make(chan data.LevelUpEvent, 10)
	levelUpResponseChan := make(chan data.LevelUpEvent, 10)

	battleUser := battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{
		UserManager: g.opts.UserManager,
		PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{
			Logger:           logger.NewChannelLogger(logChan),
			GameStateManager: g.opts.GameManager,
			LvlUpActions:     pokemon.NewLevelUpEvents(levelUpEventsChan, levelUpResponseChan),
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

		CatchPokemonEnabled: false,
		RunBattleEnabled:   false,

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

type LearnNewMoveUI struct {
	SelectedIndex   int
	ForgetMovesBtns []*widget.Clickable
}

func (g *Gui) LoadBattle(gtx layout.Context) layout.Dimensions {
	slog.Info("Starting Battle Screen")
	g.SetCurrentScreen(BattleScreen)

	go g.InitiateBattle()
	go g.ReceiveChannelActions(gtx)

	return g.RenderBattleScreen(gtx)
}

func (g *Gui) InitiateBattle() {
	err := g.Battle.PokemonBattle.Introduction()
	if err != nil {
		slog.Error("Error found within battle", "error", err)
		panic(err)
	}

	g.SetActionBtns(EndBattle)
}

func (g *Gui) ReceiveChannelActions(gtx layout.Context) {
	for {
		select {
		case log := <-g.Battle.LogChan:
			slog.Info(log)
			g.Battle.LogContent = append(g.Battle.LogContent, log)

		case events := <-g.Battle.LevelUpEventsChan:
			slog.Info("receiving level up events", "data", events)
			g.Battle.LevelUpBody = events.Body
			if events.EventType == data.LevelUpEventEvolve {
				g.Battle.DialogActionArea = EvolveDialog
			} else if events.EventType == data.LevelUpEventLearnMove {
				slog.Info("receiving level up events", "data", events, "body", events.Body.(data.EventLearnMoveBody))
				g.Battle.DialogActionArea = LearnMoveDialog
			}
		}
	}
}

func (g *Gui) SendLevelUpResponse(gtx layout.Context, response data.LevelUpEvent) {
	slog.Info("sending level up response", "data", response)
	g.Battle.LevelUpResponseChan <- response
}

func (g *Gui) SetActionBtns(actionArea ActionAreaType) {
	g.Battle.ActionArea = actionArea
}
