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
	WaitForUserInput() *data.BattleAction
}

type BattleUserImpl struct {
	UserInputChan <-chan *data.BattleAction
}

func NewBattleUser(userInputChan <-chan *data.BattleAction) BattleUser {
	return &BattleUserImpl{
		UserInputChan: userInputChan,
	}
}   

func (bu *BattleUserImpl) WaitForUserInput() *data.BattleAction {
	return <-bu.UserInputChan
}
