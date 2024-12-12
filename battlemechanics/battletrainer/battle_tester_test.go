package battletrainer_test

import (
	"os"
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battletrainer"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/logger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
    // Initialize structured logging for tests
    logger.InitLogger()

    // Run tests
    code := m.Run()
    os.Exit(code)
}

func TestBattleTester_TestNewBattleTester(t *testing.T) {
	// returns nil when empty user is passed
	assert.Nil(t, battletrainer.NewBattleTester(&data.User{}))

	// returns a valid tester
	assert.NotNil(t, battletrainer.NewBattleTester(getTestUser()))
}

func TestBattleTester_GetTrainer(t *testing.T) {
	tester := createBattleTester()
	
	assert.NotNil(t, tester.GetTrainer())
	assert.Equal(t, "John Cena", tester.GetTrainer().Name)
	assert.Equal(t, 2, len(tester.GetTrainer().Party))
}

func TestBattleTester_GetActivePokemon(t *testing.T) {
	tester := createBattleTester()
	
	assert.NotNil(t, tester.GetActivePokemon())
	assert.Equal(t, data.BasePokemonId(1), tester.GetActivePokemon().Pokemon.ID)
	assert.Equal(t, "bulbasaur", tester.GetActivePokemon().Pokemon.Name)
	assert.Equal(t, 75, tester.GetActivePokemon().Pokemon.Level)
	assert.False(t, tester.GetActivePokemon().IsFainted)
	assert.False(t, tester.GetActivePokemon().CanEvolve)
	assert.Equal(t, 144, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestGetParty(t *testing.T) {
	tester := createBattleTester()
	
	assert.NotNil(t, tester.GetParty())
	assert.Equal(t, 1, len(tester.GetParty()))
	assert.Equal(t, data.BasePokemonId(4), tester.GetParty()[0].Pokemon.ID)
	assert.Equal(t, "charmander", tester.GetParty()[0].Pokemon.Name)
	assert.Equal(t, 75, tester.GetParty()[0].Pokemon.Level)
	assert.False(t, tester.GetParty()[0].IsFainted)
	assert.False(t, tester.GetActivePokemon().CanEvolve)
	assert.Equal(t, 144, tester.GetActivePokemon().BattleHP)
}

func TestBattleTester_TestIsDefeated(t *testing.T) {
	assert.Equal(t, false, createBattleTester().IsDefeated())
}

func TestBattleTester_TestHandleSwitch(t *testing.T) {
	tester := createBattleTester()

	input := data.BattleInput{
		Type: data.Switch,
		Selected: tester.GetActivePokemon(),
		Target: &tester.GetParty()[0],
	}
	tester.HandleSwitch(input)
	assert.Nil(t, input.Move)	
}

func TestBattleTester_TestHandleSwitch_NoUnfaintedPokemonInParty(t *testing.T) {
	tester := createBattleTester()
	tester.GetParty()[0].IsFainted = true

	input := data.BattleInput{
		Type: data.Switch,
		Selected: tester.GetActivePokemon(),
		Target: &tester.GetParty()[0],
	}
	tester.HandleSwitch(input)
	assert.Nil(t, input.Move)	
}

// utils
func createBattleTester() battletrainer.BattleTrainer {
	return battletrainer.NewBattleTester(getTestUser())
}
