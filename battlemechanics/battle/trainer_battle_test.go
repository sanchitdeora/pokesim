package battle_test

import (
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TestUserPath      = "/testfiles/saved/test_user.json"
	TestUser2Path     = "/testfiles/saved/test_user_2.json"
	TestUserCopyPath  = "/testfiles/saved/test_user copy.json"
	TestUserCopy2Path = "/testfiles/saved/test_user_2 copy.json"

	TestUserSavedPath = "/testfiles/saved"
	TestUserFileName  = "test_user"
	TestUser2FileName = "test_user_2"
)

func TestMain(m *testing.M) {
	// Initialize structured logging for tests
	logger.InitTestLogger()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestTrainerBattle_TrainerBattleManager(t *testing.T) {
	assert.Nil(t, battle.TrainerBattleManager(battle.TrainerBattleOpts{}, nil, nil))
	assert.Nil(t, battle.TrainerBattleManager(battle.TrainerBattleOpts{}, battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()), nil))
	assert.NotNil(t, battle.TrainerBattleManager(
		battle.TrainerBattleOpts{},
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()),
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()),
	))
}

func TestBattle(t *testing.T) {
	gsm := gamestate.NewGameStateManager(getTestUser(), TestUserSavedPath, TestUserFileName)
	gsm2 := gamestate.NewGameStateManager(getTestUser2(), TestUserSavedPath, TestUser2FileName)

	battle := createTrainerBattleManager(gsm, gsm2)
	assert.Nil(t, battle.Introduction())

	game, err := gamestate.LoadGame(TestUserPath)
	assert.NoError(t, err)
	user := game.User

	game2, err := gamestate.LoadGame(TestUser2Path)
	assert.NoError(t, err)
	user2 := game2.User

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins+user.Stats.Losses)

	assert.Equal(t, 1, user2.Stats.Battles)
	assert.Equal(t, 1, user2.Stats.Wins+user2.Stats.Losses)

	restoreDefaultGameTestUser()
}

func getTestUser() *data.User {
	gameSave, _ := utils.ReadJsonFromFile[gamestate.GameStateSave](TestUserPath)
	return gameSave.User.ToUser()
}

func getTestUser2() *data.User {
	gameSave, _ := utils.ReadJsonFromFile[gamestate.GameStateSave](TestUser2Path)
	return gameSave.User.ToUser()
}

func createTrainerBattleManager(gsm, gsm2 gamestate.GameStateManager) battle.PokemonBattle {
	usermanager := usermanagement.NewUserManager(usermanagement.UserOpts{GameState: gsm})

	return battle.TrainerBattleManager(
		battle.TrainerBattleOpts{UserManager: usermanager}, //add mocks here
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
			UserManager:    usermanager,
			PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
		}, gsm.Get().User),
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
			UserManager:    usermanagement.NewUserManager(usermanagement.UserOpts{GameState: gsm2}),
			PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
		}, gsm.Get().User))
}

func restoreDefaultGameTestUser() {
	gameSave, _ := utils.ReadJsonFromFile[gamestate.GameStateSave](TestUserCopyPath)
	utils.WriteJsonToFile(TestUserPath, gameSave)

	gameSave, _ = utils.ReadJsonFromFile[gamestate.GameStateSave](TestUserCopy2Path)
	utils.WriteJsonToFile(TestUser2Path, gameSave)
}
