package battletrainer_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func TestBattleUser_TestNewBattleUser(t *testing.T) {
	// returns nil when empty user is passed
	assert.Nil(t, battletrainer.NewBattleUser(&data.User{}, nil))

	// returns nil when nil input channel is passed
	assert.Nil(t, battletrainer.NewBattleUser(getTestUser(), nil))

	// returns a valid user
	assert.NotNil(t, battletrainer.NewBattleUser(getTestUser(), make(chan data.BattleInput)))
}

func TestBattleUser_GetTrainer(t *testing.T) {
	user := createBattleUser()
	
	assert.NotNil(t, user.GetTrainer())
	assert.Equal(t, "John Cena", user.GetTrainer().Name)
	assert.Equal(t, 2, len(user.GetTrainer().Party))
}

func TestBattleUser_TestGetActivePokemon(t *testing.T) {
	user := createBattleUser()

	assert.NotNil(t, user.GetActivePokemon())
	assert.Equal(t, data.BasePokemonId(1), user.GetActivePokemon().Pokemon.ID)
	assert.Equal(t, "bulbasaur", user.GetActivePokemon().Pokemon.Name)
	assert.Equal(t, 75, user.GetActivePokemon().Pokemon.Level)
	assert.False(t, user.GetActivePokemon().IsFainted)
	assert.False(t, user.GetActivePokemon().CanEvolve)
	assert.Equal(t, 144, user.GetActivePokemon().BattleHP)
}

func TestBattleUser_TestGetParty(t *testing.T) {
	user := createBattleUser()

	assert.NotNil(t, user.GetParty())
	assert.Equal(t, 1, len(user.GetParty()))
	assert.Equal(t, data.BasePokemonId(4), user.GetParty()[0].Pokemon.ID)
	assert.Equal(t, "charmander", user.GetParty()[0].Pokemon.Name)
	assert.Equal(t, 75, user.GetParty()[0].Pokemon.Level)
	assert.False(t, user.GetParty()[0].IsFainted)
	assert.False(t, user.GetActivePokemon().CanEvolve)
	assert.Equal(t, 144, user.GetActivePokemon().BattleHP)
}

func TestBattleUser_TestIsDefeated(t *testing.T) {
	assert.Equal(t, false, createBattleUser().IsDefeated())
}

func TestBattleUser_TestHandleAction(t *testing.T) {
	// TODO: add test
	assert.Equal(t, 1, 1)
}

func TestBattleUser_TestSendBattleLog(t *testing.T) {
	// TODO: add test
	assert.Equal(t, 1, 1)
}

// utils
func createBattleUser() battletrainer.BattleTrainer {
	return battletrainer.NewBattleUser(getTestUser(), make(chan data.BattleInput))
}

func getTestUser() *data.User {
	user, _ := utils.ReadJsonFromFile[data.UserSave]("/testfiles/user/test_user.json")
	return user.ToUser()
}
