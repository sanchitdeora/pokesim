package battle_test

import (
	"testing"

	battle "github.com/sanchitdeora/PokeSim/battle_v1"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/trainermanagement"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/stretchr/testify/assert"
)

func createTrainerBattle() battle.BattleSequence {
	trainer := trainermanagement.NewTrainerManager(
		trainermanagement.TrainerOpts{
			SavedTrainerPath: "C:\\Projects\\Go-projects\\src\\PokéSim\\testFiles\\test_trainer.json",
		},
	)
	
	return battle.NewTrainerBattle(&battle.TrainerBattleOpts{
		UserService: usermanagement.NewUserManager(usermanagement.UserOpts{
			SavedUserPath: "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json",
		}),
		PokemonService: pokemon.NewPokemonService(pokemon.PokemonOpts{}),
	}, trainer.GetTrainer())
}

func TestBattle(t *testing.T) {
	battle := createTrainerBattle()

	report, err := battle.Initiate()
	assert.NoError(t, err)
	if report.UserWin {
		assert.Equal(t, 6000, report.Money)
	} else {
		assert.Equal(t, 600, report.Money)
	}
}
