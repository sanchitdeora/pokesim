package battletrainer

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
)

// user
type BattleUser struct {
	BattleTrainerImpl

	actionChan <-chan data.BattleAction
}

func NewBattleUser(opts BattleTrainerOpts, user *data.User, actionChan <-chan data.BattleAction, logsChan chan<- string) BattleTrainer {
	if len(user.Party) == 0 {
		slog.Error("user party is nil")
		return nil
	}
	if actionChan == nil {
		slog.Error("action channel is nil")
		return nil
	}

	party := make([]*data.BattlePokemon, 0)
	for _, p := range user.Party[1:] {
		party = append(party, data.CreateBattlePokemon(p))
	}

	return &BattleUser{
		BattleTrainerImpl: BattleTrainerImpl{
			BattleTrainerOpts: opts,
			Trainer:           user,
			ActivePokemon:     data.CreateBattlePokemon(user.Party[0]),
			BattleParty:       party,
			Rewards:           data.Rewards{},
			Logger:            logger.NewChannelLogger(logsChan),
		},
		actionChan: actionChan,
	}
}

func (b *BattleUser) HandleBattleEnd(result data.Result) {
	b.handleBattleEnd(result)
}

func (b *BattleUser) GetAction() data.BattleAction {
	return <-b.actionChan
}
