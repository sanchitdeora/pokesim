package gamestate_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/gamestate"
	"github.com/stretchr/testify/assert"
)

func createGameState() gamestate.GameStateManager {
	return gamestate.NewGameStateManager(nil, "testfiles/gametest", "unittestgame")
}

func TestNewGamestate_WithNilConstraints(t *testing.T) {
	gameManager := gamestate.NewGameStateManager(nil, "testfiles/gametest", "nilUserTest")
	assert.NotNil(t, gameManager)
	assert.Equal(t, gameManager.Get().User.Name, "")

	gameManager = gamestate.NewGameStateManager(&data.User{}, "testfiles/gametest", "defaultUserTest")
	assert.NotNil(t, gameManager)
	assert.Equal(t, gameManager.Get().User.Name, "")
}

func TestNewGamestate_Happy(t *testing.T) {
	gameManager := gamestate.NewGameStateManager(&data.User{BaseTrainer: data.BaseTrainer{Name: "test"}}, "", "")
	assert.NotNil(t, gameManager)
	assert.Equal(t, gameManager.Get().User.Name, "test")

	gameManager = gamestate.NewGameStateManager(&data.User{BaseTrainer: data.BaseTrainer{Name: "test"}}, "testfiles/gametest", "ash1")
	assert.NotNil(t, gameManager)
	assert.Equal(t, gameManager.Get().User.Name, "test")
}

// func TestAddTrainerProgress(t *testing.T) {
// 	gameManager := createGameState()

// 	gameManager.AddTrainerProgress("test1Trainer")

// 	game, err := gameManager.Load()
// 	assert.NoError(t, err)
// 	assert.Equal(t, game.TrainerProgress[len(game.TrainerProgress)-1], "test1Trainer")
// }
