package battle

import "github.com/sanchitdeora/PokeSim/data"

type BattleSequence interface {
	Initiate() (*data.Result, error)
	Attack(attackPokemon *data.BattlePokemon, targetPokemon *data.BattlePokemon, attackMove *data.Moves, isUser bool)
	CatchPokemon(targetPokemon *data.BattlePokemon, item *data.Item)
	Run()
	IsOver() bool
	Report() (*data.Result, error)
	GetUserActive(isUser bool) *data.BattlePokemon
}

type BattleUser interface {
	WaitForUserInput() *data.BattleInput
}

type BattleUserImpl struct {
	UserInputChan <-chan *data.BattleInput
}

func NewBattleUser(userInputChan <-chan *data.BattleInput) BattleUser {
	return &BattleUserImpl{
		UserInputChan: userInputChan,
	}
}   

func (bu *BattleUserImpl) WaitForUserInput() *data.BattleInput {
	return <-bu.UserInputChan
}
