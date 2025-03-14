package battletrainer

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
)

type BattleWild struct {
	BattleTrainerImpl
}

func NewBattleWild(opts BattleTrainerOpts, wild *data.Pokemon) BattleTrainer {
	if wild == nil {
		slog.Error("wild pokemon is nil")
		return nil
	}

	wildTrainer := &data.BaseTrainer{
		Name:  "Wild Trainer",
		Party: []*data.Pokemon{wild},
		Bag:   data.ItemMap{},
	}

	return &BattleWild{
		BattleTrainerImpl: BattleTrainerImpl{
			BattleTrainerOpts: opts,
			Trainer:           wildTrainer,
			ActivePokemon:     data.CreateBattlePokemon(wild),
			BattleParty:       make([]*data.BattlePokemon, 0),
			Rewards:           data.Rewards{},
			Logger:            logger.NewDefaultLogger(),
		},
	}
}

func (b *BattleWild) GetTrainerType() data.TrainerClass {
	return data.WildPrefix
}

func (b *BattleWild) IsDefeated() bool {
	return b.ActivePokemon.BattleHP == 0
}


func (b *BattleWild) HandleTargetPokemonFainted(faintedPokemon *data.BattlePokemon) {
	return
}

func (b *BattleWild) HandleUseBag(action data.BattleAction) error {
	return nil
}

func (b *BattleWild) GetAction() data.BattleAction {
	slog.Info("Battle Wild getting action")

	return data.BattleAction{
		ID:       uuid.NewString(),
		Move:     randomMove(&b.GetActivePokemon().Moveset),
		Type:     data.Attack,
		Selected: b.GetActivePokemon(),
		Target:   b.Opponent,
	}
}

func (b *BattleWild) CalculateResult(opponent BattleTrainer) data.Result {
	return data.Result{}
}

func (b *BattleWild) HandleBattleEnd(result data.Result) {
	slog.Info("Battle AI handling end")
}
