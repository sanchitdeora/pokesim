package battletrainer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	mock_pokemon_manager "github.com/sanchitdeora/PokeSim/pokemon/mocks"
	mock_user_manager "github.com/sanchitdeora/PokeSim/usermanagement/mocks"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func TestBattleUser_TestNewBattleUser(t *testing.T) {
	logChan := make(chan<- string)
	actionChan := make(chan data.BattleAction)
	// returns nil when empty user is passed
	assert.Nil(t, battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{}, &data.User{}, nil, logChan))

	// returns nil when nil action channel is passed
	assert.Nil(t, battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{}, getTestUser(), nil, logChan))

	// returns a valid user
	assert.NotNil(t, battletrainer.NewBattleUser(battletrainer.BattleTrainerOpts{}, getTestUser(),actionChan, logChan))
	
	close(actionChan)
	close(logChan)
}

func TestBattleUser_GetTrainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	user := createBattleUser(ctrl)

	assert.NotNil(t, user.GetTrainer())
	assert.Equal(t, "John Cena", user.GetTrainer().Name)
	assert.Equal(t, 2, len(user.GetTrainer().Party))
}

func TestBattleUser_TestGetActivePokemon(t *testing.T) {
	ctrl := gomock.NewController(t)
	user := createBattleUser(ctrl)

	assert.NotNil(t, user.GetActivePokemon())
	assert.Equal(t, data.BasePokemonID(1), user.GetActivePokemon().Pokemon.ID)
	assert.Equal(t, "bulbasaur", user.GetActivePokemon().Pokemon.Name)
	assert.Equal(t, 75, user.GetActivePokemon().Pokemon.Level)
	assert.False(t, user.GetActivePokemon().CanEvolve)
	assert.Equal(t, 152, user.GetActivePokemon().BattleHP)
}

func TestBattleUser_TestGetParty(t *testing.T) {
	ctrl := gomock.NewController(t)
	user := createBattleUser(ctrl)

	assert.NotNil(t, user.GetParty())
	assert.Equal(t, 1, len(user.GetParty()))
	assert.Equal(t, data.BasePokemonID(4), user.GetParty()[0].Pokemon.ID)
	assert.Equal(t, "charmander", user.GetParty()[0].Pokemon.Name)
	assert.Equal(t, 75, user.GetParty()[0].Pokemon.Level)
	assert.False(t, user.GetActivePokemon().CanEvolve)
	assert.Equal(t, 152, user.GetActivePokemon().BattleHP)
}

func TestBattleUser_TestIsDefeated(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert.Equal(t, false, createBattleUser(ctrl).IsDefeated())
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
func createBattleUser(ctrl *gomock.Controller) battletrainer.BattleTrainer {
	return battletrainer.NewBattleUser(
		battletrainer.BattleTrainerOpts{
			PokemonService: mock_pokemon_manager.NewMockPokemonService(ctrl),
			UserManager:    mock_user_manager.NewMockUserManager(ctrl),
		},
		getTestUser(),
		make(chan data.BattleAction),
		make(chan<- string),
	)
}

func getTestUser() *data.User {
	user, _ := utils.ReadJsonFromFile[data.UserSave]("/testfiles/user/test_user.json")
	return user.ToUser()
}
