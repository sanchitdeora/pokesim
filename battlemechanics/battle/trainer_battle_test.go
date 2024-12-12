package battle_test

import (
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize structured logging for tests
	logger.InitLogger()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestTrainerBattle_TrainerBattleManager(t *testing.T) {
	assert.Nil(t, battle.TrainerBattleManager(battle.TrainerBattleOpts{}, nil, nil))
	assert.Nil(t, battle.TrainerBattleManager(battle.TrainerBattleOpts{}, battletrainer.NewBattleTester(getTestUser()), nil))
	assert.NotNil(t, battle.TrainerBattleManager(battle.TrainerBattleOpts{}, battletrainer.NewBattleTester(getTestUser()), battletrainer.NewBattleTester(getTestUser())))
}

func TestIntroduction(t *testing.T) {
	battle := createTrainerBattleManager()
	assert.NotNil(t, battle.Introduction())
}

func getTestUser() *data.User {
	user, _ := utils.ReadJsonFromFile[data.UserSave]("/testfiles/user/test_user.json")
	return user.ToUser()
}

func createTrainerBattleManager() battle.PokemonBattle {
	return battle.TrainerBattleManager(
		battle.TrainerBattleOpts{UserManager: nil}, //add mocks here
		battletrainer.NewBattleTester(getTestUser()),
		battletrainer.NewBattleTester(getTestUser()))
}
