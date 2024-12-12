package battletrainer

import "github.com/sanchitdeora/PokeSim/data"

// type battletrainer interface {
// 	WaitForInput(pokemon *data.BattlePokemon, target *data.BattlePokemon) *data.BattleInput
// }

type BattleTrainer interface {
	GetTrainer() data.BaseTrainer
	GetActivePokemon() *data.BattlePokemon
	GetParty() []data.BattlePokemon	
	GetTrainerType() data.TrainerClass
	GetRewards() data.Rewards
	
	IsDefeated() bool
	GetInput() data.BattleInput
	HandleSwitch(action data.BattleInput) error
	HandleUseBag(action data.BattleInput) error

	SendBattleLog(message string) error

	// Only useful for testing
	SetOpponentTarget(target *data.BattlePokemon)
}
