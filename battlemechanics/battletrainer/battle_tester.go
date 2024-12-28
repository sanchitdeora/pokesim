package battletrainer

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
)

type BattleTester struct {
	BattleTrainerImpl
}

func NewBattleTester(opts BattleTrainerOpts, user *data.User) BattleTrainer {
	if len(user.Party) == 0 {
		slog.Error("user party is nil")
		return nil
	}

	party := make([]*data.BattlePokemon, 0)
	for _, p := range user.Party[1:] {
		party = append(party, data.CreateBattlePokemon(p))
	}

	return &BattleTester{
		BattleTrainerImpl: BattleTrainerImpl{
			BattleTrainerOpts: opts,
			Trainer:           user,
			ActivePokemon:     data.CreateBattlePokemon(user.Party[0]),
			BattleParty:       party,
			Rewards:           data.Rewards{},
			Logger:            logger.NewDefaultLogger(),
		},
	}
}

func (b *BattleTester) GetAction() data.BattleAction {
	return data.BattleAction{
		ID:       uuid.NewString(),
		Move:     randomMove(&b.GetActivePokemon().Moveset),
		Type:     data.Attack,
		Selected: b.GetActivePokemon(),
		Target:   b.Opponent,
	}
}
