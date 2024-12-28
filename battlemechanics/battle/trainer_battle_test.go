package battle_test

import (
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TestUserPath      = "/testfiles/user/test_user.json"
	TestUser2Path     = "/testfiles/user/test_user_2.json"
	TestUserCopyPath  = "/testfiles/user/test_user copy.json"
	TestUserCopy2Path = "/testfiles/user/test_user_2 copy.json"
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
	battle := createTrainerBattleManager()
	assert.Nil(t, battle.Introduction())

	user := getTestUser()
	user2 := getTestUser2()

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins+user.Stats.Losses)

	assert.Equal(t, 1, user2.Stats.Battles)
	assert.Equal(t, 1, user2.Stats.Wins+user2.Stats.Losses)

	restoreDefaultTestUser()
}

func getTestUser() *data.User {
	user, _ := utils.ReadJsonFromFile[data.UserSave](TestUserPath)
	return user.ToUser()
}

func getTestUser2() *data.User {
	user, _ := utils.ReadJsonFromFile[data.UserSave](TestUser2Path)
	return user.ToUser()
}

func createTrainerBattleManager() battle.PokemonBattle {
	usermanager := usermanagement.NewUserService(usermanagement.UserOpts{SavedUserPath: TestUserPath})

	return battle.TrainerBattleManager(
		battle.TrainerBattleOpts{UserManager: usermanager}, //add mocks here
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
			UserManager:    usermanager,
			PokemonManager: pokemon.NewPokemonManager(pokemon.PokemonOpts{}),
		}, getTestUser()),
		battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
			UserManager:    usermanagement.NewUserService(usermanagement.UserOpts{SavedUserPath: TestUser2Path}),
			PokemonManager: pokemon.NewPokemonManager(pokemon.PokemonOpts{}),
		}, getTestUser2()))
}

func restoreDefaultTestUser() {
	userSave, _ := utils.ReadJsonFromFile[data.UserSave](TestUserCopyPath)
	utils.WriteJsonToFile(TestUserPath, userSave)

	userSave, _ = utils.ReadJsonFromFile[data.UserSave](TestUserCopy2Path)
	utils.WriteJsonToFile(TestUser2Path, userSave)
}
