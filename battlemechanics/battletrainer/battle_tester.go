package battletrainer

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/sanchitdeora/PokeSim/data"
)

type BattleTester struct {
	Trainer       data.BaseTrainer
	ActivePokemon *data.BattlePokemon
	BattleParty   []data.BattlePokemon
	Opponent      *data.BattlePokemon
}

func NewBattleTester(user *data.User) BattleTrainer {
	if len(user.Party) == 0 {
		slog.Error("user party is nil")
		return nil
	}

	party := make([]data.BattlePokemon, 0)
	for _, p := range user.Party[1:] {
		party = append(party, *createBattlePokemon(p))
	}

	return &BattleTester{
		Trainer:       user.BaseTrainer,
		ActivePokemon: createBattlePokemon(user.Party[0]),
		BattleParty:   party,
	}
}

func (b *BattleTester) GetTrainer() data.BaseTrainer {
	return b.Trainer
}

func (b *BattleTester) GetActivePokemon() *data.BattlePokemon {
	return b.ActivePokemon
}

func (b *BattleTester) GetParty() []data.BattlePokemon {
	return b.BattleParty
}

func (b *BattleTester) GetTrainerType() data.TrainerClass {
	return data.TrainerPrefix
}

func (b *BattleTester) GetRewards() data.Rewards {
	return data.Rewards{}
}

func (b *BattleTester) GetInput() data.BattleInput {
	return data.BattleInput{
		Move:     randomMove(&b.GetActivePokemon().Moveset),
		Type:     data.Attack,
		Selected: b.GetActivePokemon(),
		Target:   b.Opponent,
	}
}

func (b *BattleTester) HandleSwitch(action data.BattleInput) error {
	b.ActivePokemon, b.BattleParty = switchPokemon(b.ActivePokemon, action.Target, b.BattleParty)

	slog.Info(fmt.Sprintf("%s is switching their pokemon!", b.Trainer.Name))
	slog.Info(fmt.Sprintf("%s, chooses %s!", b.Trainer.Name, action.Selected.Pokemon.Name))

	return nil
}

func (b *BattleTester) HandleUseBag(action data.BattleInput) error {
	slog.Info("Handling use bag", "action", action.Type)
	// TODO: Implement logic for handling attack, switch, item use, etc.
	return nil
}

func (b *BattleTester) IsDefeated() bool {
	return len(b.GetParty()) == 0
}

func (b *BattleTester) SendBattleLog(message string) error {
	slog.Debug(message)
	return nil
}

func (b *BattleTester) SetOpponentTarget(target *data.BattlePokemon) {
	b.Opponent = target
}

func randomMove(moveset *data.Moveset) *data.Moves {
	moves := []*data.Moves{}
	if moveset.Move1 != nil {
		moves = append(moves, moveset.Move1)
	}
	if moveset.Move2 != nil {
		moves = append(moves, moveset.Move2)
	}
	if moveset.Move3 != nil {
		moves = append(moves, moveset.Move3)
	}
	if moveset.Move4 != nil {
		moves = append(moves, moveset.Move4)
	}

	if len(moves) == 0 {
		return nil
	}

	input := randomGenerator(0, float64(len(moves)))
	return moves[int(input)]
}

func randomGenerator(min float64, max float64) float64 {
	randIndex := rand.Float64()
	randomGenerator := (min + randIndex*(max-min))
	if randomGenerator == max {
		return randomGenerator - 1
	}
	return randomGenerator
}
