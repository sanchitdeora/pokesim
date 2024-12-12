package battletrainer

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

// user
type BattleUser struct {
	Trainer       *data.User
	ActivePokemon *data.BattlePokemon
	BattleParty   []data.BattlePokemon

	inputChan <-chan data.BattleInput
	LogsChan  chan<- string // Sends logs back to UI

	Opponent *data.BattlePokemon
}

func NewBattleUser(user *data.User, inputChan <-chan data.BattleInput) BattleTrainer {
	if len(user.Party) == 0 {
		slog.Error("user party is nil")
		return nil
	}
	if inputChan == nil {
		slog.Error("input channel is nil")
		return nil
	}

	party := make([]data.BattlePokemon, 0)
	for _, p := range user.Party[1:] {
		party = append(party, *createBattlePokemon(p))
	}

	return &BattleUser{
		Trainer:       user,
		ActivePokemon: createBattlePokemon(user.Party[0]),
		BattleParty:   party,
		inputChan:     inputChan,
	}
}

func (b *BattleUser) GetTrainer() data.BaseTrainer {
	return b.Trainer.BaseTrainer	
}

func (bu *BattleUser) GetTrainerType() data.TrainerClass {
	return data.RivalPrefix
}

func (bu *BattleUser) GetActivePokemon() *data.BattlePokemon {
	return bu.ActivePokemon
}

func (bu *BattleUser) GetParty() []data.BattlePokemon {
	return bu.BattleParty
}

func (bu *BattleUser) GetRewards() data.Rewards {
	return data.Rewards{}
}

func (b *BattleUser) IsDefeated() bool {
	return len(b.GetParty()) == 0
}

func (b *BattleUser) HandleSwitch(action data.BattleInput) error {
	b.ActivePokemon, b.BattleParty = switchPokemon(b.ActivePokemon, action.Target, b.BattleParty)

	b.SendBattleLog(fmt.Sprintf("%s is switching their pokemon!", b.Trainer.Name))
	b.SendBattleLog(fmt.Sprintf("%s, chooses %s!", b.Trainer.Name, action.Selected.Pokemon.Name))

	return nil
}

func (b *BattleUser) HandleUseBag(action data.BattleInput) error {
	slog.Info("Handling use bag", "action", action.Type)
	// TODO: Implement logic for handling attack, switch, item use, etc.
	return nil
}

func (b *BattleUser) HandleBattleEnd() error {
	panic("add user stats etc.")
}

func (b *BattleUser) SendBattleLog(message string) error {
	return nil
}

func (bu *BattleUser) GetInput() data.BattleInput {
	return <-bu.inputChan
}

func (bu *BattleUser) SetOpponentTarget(target *data.BattlePokemon) {
	bu.Opponent = target
}
