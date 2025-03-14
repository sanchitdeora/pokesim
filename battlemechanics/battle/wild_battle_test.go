package battle_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/stretchr/testify/assert"
)

const (
)

func TestWildBattle_WildBattleManager(t *testing.T) {
	assert.Nil(t, battle.WildBattleManager(battle.WildBattleOpts{}, nil, nil))
	assert.Nil(t, battle.WildBattleManager(battle.WildBattleOpts{}, battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()), nil))
	assert.NotNil(t, battle.WildBattleManager(
		battle.WildBattleOpts{},
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()),
		battletrainer.NewBattleWild(battletrainer.BattleTrainerOpts{}, getTestUser().Party[0]),
	))

	assert.Nil(t, battle.WildBattleManager(
		battle.WildBattleOpts{},
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()),
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()),
	))
}

func TestWildBattle(t *testing.T) {
	gsm := gamestate.NewGameStateManager(getTestUser(), TestUserSavedPath, TestUserFileName)

	battle := createWildBattleManager(gsm)
	assert.Nil(t, battle.Introduction())

	game, err := gamestate.LoadGame(TestUserPath)
	assert.NoError(t, err)
	user := game.User

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins+user.Stats.Losses)

	restoreDefaultGameTestUser()
}

func createWildBattleManager(gsm gamestate.GameStateManager) battle.PokemonBattle {
	usermanager := usermanagement.NewUserManager(usermanagement.UserOpts{GameState: gsm})

	return battle.WildBattleManager(
		battle.WildBattleOpts{UserManager: usermanager}, //add mocks here
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
			UserManager:    usermanager,
			PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
		}, gsm.Get().User),
		battletrainer.NewBattleWild(battletrainer.BattleTrainerOpts{
			UserManager:    nil,
			PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
		}, gsm.Get().User.Party[0]))
}
