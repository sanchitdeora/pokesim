package battletrainer

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
)

type BattleAI struct {
	BattleTrainerImpl
}

func NewBattleAI(opts BattleTrainerOpts, trainer *data.Trainer) BattleTrainer {
	if len(trainer.Party) == 0 {
		slog.Error("user party is nil")
		return nil
	}

	party := make([]*data.BattlePokemon, 0)
	for _, p := range trainer.Party[1:] {
		party = append(party, data.CreateBattlePokemon(p))
	}

	return &BattleAI{
		BattleTrainerImpl: BattleTrainerImpl{
			BattleTrainerOpts: opts,
			Trainer:           &trainer.BaseTrainer,
			ActivePokemon:     data.CreateBattlePokemon(trainer.Party[0]),
			BattleParty:       party,
			Rewards:           *trainer.Rewards,
			Logger:            logger.NewDefaultLogger(),
		},
	}
}

func (b *BattleAI) HandleUseBag(action data.BattleAction) error {
	slog.Info("Healing pokemon", "pokemon", action.Target.PokemonUUID)
	var target *data.BattlePokemon
	if action.Target.PokemonUUID == b.ActivePokemon.PokemonUUID {
		target = b.ActivePokemon
	} else {
		for _, p := range b.BattleParty {
			if p.PokemonUUID == action.Target.PokemonUUID {
				target = p
				break
			}
		}
	}

	// TODO: implement healing for trainer management
	if healPokemon(target, action.Item) {
		b.UserManager.UseItem(action.Item)
	}
	return nil
}

func (b *BattleAI) HandleTargetPokemonFainted(faintedPokemon *data.BattlePokemon) {
	// do nothing
}

func (b *BattleAI) GetAction() data.BattleAction {
	slog.Info("Battle AI getting action")

	return data.BattleAction{
		ID:       uuid.NewString(),
		Move:     randomMove(&b.GetActivePokemon().Moveset),
		Type:     data.Attack,
		Selected: b.GetActivePokemon(),
		Target:   b.Opponent,
	}
}

func (b *BattleAI) CalculateResult(opponent BattleTrainer) data.Result {
	return data.Result{}
}

func (b *BattleAI) HandleBattleEnd(result data.Result) {
	slog.Info("Battle AI handling end")
}
