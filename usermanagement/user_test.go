package usermanagement_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	mock_game_state_manager "github.com/sanchitdeora/PokeSim/gamestate/mocks"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TestUserPath     = "testfiles\\saved\\unit_test_user.json"
	TestUserCopyPath = "testfiles\\saved\\unit_test_user copy.json"
)

func createUserManager(mockGameState *mock_game_state_manager.MockGameStateManager) usermanagement.UserManager {
	return usermanagement.NewUserManager(usermanagement.UserOpts{GameState: mockGameState})
}

// LoadUser
func TestLoadUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)

	user := userManager.GetUser()

	assert.Equal(t, "John Cena", user.Name)
	assert.Equal(t, "bulbasaur", user.Party[0].Name)
	assert.Equal(t, "charmander", user.Party[1].Name)
}

// StatUpdate
func TestStatUpdate_BattleWon(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)

	user := userManager.GetUser()

	result := data.Result{
		Status: data.Won,
		Money:  100,
		BonusItems: data.ItemMap{
			data.Potion:   {Category: data.MedicalItems, Count: 1},
			data.PokeBall: {Category: data.PokeBalls, Count: 1},
		},
		BadgeEarned: data.BadgeType{
			Name:   "test",
			Region: "test",
		},
	}

	gsm.EXPECT().Save().Return(nil)

	userManager.StatUpdate(result)

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins)
	assert.Equal(t, 0, user.Stats.Losses)
	assert.Equal(t, 1, len(user.Stats.Badges))
	assert.Equal(t, 100, user.Money)
	assert.Equal(t, 2, len(user.Bag))
	assert.Equal(t, 3, user.Bag[data.Potion].Count)
	assert.Equal(t, 1, user.Bag[data.PokeBall].Count)

	// clean up after
	restoreDefaultGameTestUser()
}

func TestStatUpdate_BattleLost_NoMoney(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)

	user := userManager.GetUser()

	result := data.Result{
		Status: data.Lost,
		Money:  100,
	}

	gsm.EXPECT().Save().Return(nil)

	userManager.StatUpdate(result)

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 0, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Losses)
	assert.Equal(t, 0, len(user.Stats.Badges))
	assert.Equal(t, 0, user.Money)
	assert.Equal(t, 1, len(user.Bag))
	assert.Equal(t, 2, user.Bag[data.Potion].Count)

	// clean up after
	restoreDefaultGameTestUser()
}

func TestStatUpdate_BattleLost_LostMoney(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)

	user := userManager.GetUser()
	user.Money = 500

	result := data.Result{
		Status: data.Lost,
		Money:  100,
	}

	gsm.EXPECT().Save().Return(nil)

	userManager.StatUpdate(result)

	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 0, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Losses)
	assert.Equal(t, 0, len(user.Stats.Badges))
	assert.Equal(t, 400, user.Money)
	assert.Equal(t, 1, len(user.Bag))
	assert.Equal(t, 2, user.Bag[data.Potion].Count)

	// clean up after
	restoreDefaultGameTestUser()
}

// ChangePokemonOrder
func TestChangePokemonOrder_MoveDown(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)
	user := userManager.GetUser()

	firstPokemon := user.Party[0]
	secondPokemon := user.Party[1]

	gsm.EXPECT().Save().Return(nil)

	userManager.ChangePokemonOrder(firstPokemon, data.ChangeOrderMoveDown)

	assert.Equal(t, secondPokemon, user.Party[0])
	assert.Equal(t, firstPokemon, user.Party[1])
}

func TestChangePokemonOrder_MoveUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)
	user := userManager.GetUser()

	firstPokemon := user.Party[0]
	secondPokemon := user.Party[1]

	gsm.EXPECT().Save().Return(nil)

	userManager.ChangePokemonOrder(secondPokemon, data.ChangeOrderMoveUp)

	assert.Equal(t, secondPokemon, user.Party[0])
	assert.Equal(t, firstPokemon, user.Party[1])
}

func TestChangePokemonOrder_MoveDownWthLastPokemon_NoChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)
	user := userManager.GetUser()

	firstPokemon := user.Party[0]
	secondPokemon := user.Party[1]

	gsm.EXPECT().Save().Return(nil)

	userManager.ChangePokemonOrder(secondPokemon, data.ChangeOrderMoveDown)

	assert.Equal(t, firstPokemon, user.Party[0])
	assert.Equal(t, secondPokemon, user.Party[1])
}

func TestChangePokemonOrder_MoveUpWthFirstPokemon_NoChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	gsm := mock_game_state_manager.NewMockGameStateManager(ctrl)

	gsm.EXPECT().Get().Return(getTestGameState())

	userManager := createUserManager(gsm)
	user := userManager.GetUser()

	firstPokemon := user.Party[0]
	secondPokemon := user.Party[1]

	gsm.EXPECT().Save().Return(nil)

	userManager.ChangePokemonOrder(firstPokemon, data.ChangeOrderMoveUp)

	assert.Equal(t, firstPokemon, user.Party[0])
	assert.Equal(t, secondPokemon, user.Party[1])
}

func restoreDefaultGameTestUser() {
	gameSave, _ := utils.ReadJsonFromFile[gamestate.GameStateSave](TestUserCopyPath)
	utils.WriteJsonToFile(TestUserPath, gameSave)
}

func getTestGameState() *gamestate.GameState {
	game, _ := utils.ReadJsonFromFile[gamestate.GameStateSave](TestUserPath)
	return game.ToGameState()
}
