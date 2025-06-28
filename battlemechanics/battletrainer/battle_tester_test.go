package battletrainer_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	mock_pokemon_manager "github.com/sanchitdeora/PokeSim/pokemon/mocks"
	mock_user_manager "github.com/sanchitdeora/PokeSim/usermanagement/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize structured logging for tests
	logger.InitTestLogger()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestBattleTester_TestNewBattleTester(t *testing.T) {
	// returns nil when empty user is passed
	assert.Nil(t, battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, &data.User{}))

	// returns a valid tester
	assert.NotNil(t, battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{}, getTestUser()))
}

func TestBattleTester_GetTrainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	assert.NotNil(t, tester.GetTrainer())
	assert.Equal(t, "John Cena", tester.GetTrainer().Name)
	assert.Equal(t, 2, len(tester.GetTrainer().Party))
}

func TestBattleTester_GetActivePokemon(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	assert.NotNil(t, tester.GetActivePokemon())
	assert.Equal(t, data.BasePokemonID(1), tester.GetActivePokemon().Pokemon.ID)
	assert.Equal(t, "bulbasaur", tester.GetActivePokemon().Pokemon.Name)
	assert.Equal(t, 75, tester.GetActivePokemon().Pokemon.Level)
	assert.False(t, tester.GetActivePokemon().IsFainted)
	assert.False(t, tester.GetActivePokemon().CanEvolve)
	assert.Equal(t, 152, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestGetParty(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	assert.NotNil(t, tester.GetParty())
	assert.Equal(t, 1, len(tester.GetParty()))
	assert.Equal(t, data.BasePokemonID(4), tester.GetParty()[0].Pokemon.ID)
	assert.Equal(t, "charmander", tester.GetParty()[0].Pokemon.Name)
	assert.Equal(t, 75, tester.GetParty()[0].Pokemon.Level)
	assert.False(t, tester.GetParty()[0].IsFainted)
	assert.False(t, tester.GetActivePokemon().CanEvolve)
	assert.Equal(t, 152, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestIsDefeated(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	assert.False(t, tester.IsDefeated())

	for i := range tester.GetParty() {
		tester.GetParty()[i].BattleHP = 0
	}
	assert.False(t, tester.IsDefeated())

	tester.GetActivePokemon().BattleHP = 0
	assert.True(t, tester.IsDefeated())
}

func TestBattleTester_TestHandleSwitch(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	action := data.BattleAction{
		Type:     data.Switch,
		Selected: tester.GetActivePokemon(),
		Target:   tester.GetParty()[0],
	}
	err := tester.HandleAction(action)
	assert.Nil(t, err)
	assert.Equal(t, action.Target, tester.GetActivePokemon())
	assert.Equal(t, 1, len(tester.GetParty()))
	assert.Equal(t, action.Selected, tester.GetParty()[0])
}

func TestBattleTester_TestHandleSwitch_NoUnfaintedPokemonInParty(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	tester.GetParty()[0].BattleHP = 0

	action := data.BattleAction{
		Type:     data.Switch,
		Selected: tester.GetActivePokemon(),
		Target:   tester.GetParty()[0],
	}
	err := tester.HandleAction(action)
	assert.NotNil(t, err)
	assert.Equal(t, action.Selected, tester.GetActivePokemon())
	assert.Equal(t, 1, len(tester.GetParty()))
	assert.Equal(t, action.Target, tester.GetParty()[0])
}

func TestBattleTester_TestHandleUseBag_AlreadyFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, _ := createBattleTester(ctrl)

	item := getTestUser().Bag[data.Potion]

	expBattleHP := tester.GetActivePokemon().BattleHP

	action := data.BattleAction{
		Type:     data.Bag,
		Selected: tester.GetActivePokemon(),
		Target:   tester.GetActivePokemon(),
		Item:     &item,
	}
	err := tester.HandleAction(action)
	assert.Nil(t, err)

	assert.Equal(t, expBattleHP, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestHandleUseBag_LessThan10(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, mocks := createBattleTester(ctrl)

	item := getTestUser().Bag[data.Potion]

	expBattleHP := tester.GetActivePokemon().BattleHP
	tester.GetActivePokemon().BattleHP -= 10

	action := data.BattleAction{
		Type:     data.Bag,
		Selected: tester.GetActivePokemon(),
		Target:   tester.GetActivePokemon(),
		Item:     &item,
	}

	mocks.UserManager.EXPECT().UseItem(&item, 1).Times(1)
	err := tester.HandleAction(action)
	assert.Nil(t, err)

	assert.Equal(t, expBattleHP, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestHandleUseBag_LessThan50(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, mocks := createBattleTester(ctrl)

	item := getTestUser().Bag[data.Potion]

	tester.GetActivePokemon().BattleHP -= 50
	expBattleHP := tester.GetActivePokemon().BattleHP + 20

	action := data.BattleAction{
		Type:     data.Bag,
		Selected: tester.GetActivePokemon(),
		Target:   tester.GetActivePokemon(),
		Item:     &item,
	}

	mocks.UserManager.EXPECT().UseItem(&item, 1).Times(1)
	err := tester.HandleAction(action)
	assert.Nil(t, err)

	assert.Equal(t, expBattleHP, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestHandleUseBag_LessThan50_InParty(t *testing.T) {
	ctrl := gomock.NewController(t)
	tester, mocks := createBattleTester(ctrl)

	item := getTestUser().Bag[data.Potion]

	tester.GetParty()[0].BattleHP -= 50
	expBattleHP := tester.GetParty()[0].BattleHP + 20

	action := data.BattleAction{
		Type:     data.Bag,
		Selected: tester.GetParty()[0],
		Target:   tester.GetParty()[0],
		Item:     &item,
	}

	mocks.UserManager.EXPECT().UseItem(&item, 1).Times(1)
	err := tester.HandleAction(action)
	assert.Nil(t, err)

	assert.Equal(t, expBattleHP, tester.GetParty()[0].BattleHP)
}

// utils
func createBattleTester(ctrl *gomock.Controller) (battletrainer.BattleTrainer, MocksImpl) {
	ps := mock_pokemon_manager.NewMockPokemonService(ctrl)
	um := mock_user_manager.NewMockUserManager(ctrl)

	mocks := MocksImpl{
		PokemonService: *ps,
		UserManager:    *um,
	}

	return battletrainer.NewBattleTester(battletrainer.BattleTrainerOpts{
		PokemonService: ps,
		UserManager:    um,
	}, getTestUser()), mocks
}

type MocksImpl struct {
	PokemonService mock_pokemon_manager.MockPokemonService
	UserManager    mock_user_manager.MockUserManager
}
