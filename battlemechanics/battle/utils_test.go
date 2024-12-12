package battle_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetTurnOrder_BothSwitch(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Switch,
		Selected: GetPokemonFromId(0),
	}
	opponentInput := data.BattleInput{
		Type: data.Switch,
		Selected: GetPokemonFromId(1),
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, userInput, opponentInput, order)
}

func TestGetTurnOrder_BothUseBag(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Bag,
		Selected: GetPokemonFromId(0),
	}
	opponentInput := data.BattleInput{
		Type: data.Bag,
		Selected: GetPokemonFromId(1),
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, userInput, opponentInput, order)
}

func TestGetTurnOrder_BothAttack_UserHighPriorityMove(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(1),
		Move: &data.Moves{
			Priority: -1,
		},
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, userInput, opponentInput, order)
}

func TestGetTurnOrder_BothAttack_OpponentHighPriorityMove(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(0),
		Move: &data.Moves{
			Priority: -1,
		},
	}
	opponentInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, opponentInput, userInput, order)
}

func TestGetTurnOrder_BothAttack_UserHighSpeed(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, userInput, opponentInput, order)
}

func TestGetTurnOrder_BothAttack_OpponentHighSpeed(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(1),
		Move: &data.Moves{
			Name: "Flamethrower",
			Priority: 1,
		},
	}
	opponentInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(0),
		Move: &data.Moves{
			Name: "Water Gun",
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, opponentInput, userInput, order)
}

func TestGetTurnOrder_UserAttackOpponentSwitch_UserHighSpeed(t *testing.T) {
	userInput := data.BattleInput{
		Type: data.Attack,
		Selected: GetPokemonFromId(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentInput := data.BattleInput{
		Type: data.Switch,
		Selected: GetPokemonFromId(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userInput, opponentInput)

	assertTurnOrder(t, opponentInput, userInput, order)
}

func assertTurnOrder(t *testing.T, expFirstInput data.BattleInput, expSecondInput data.BattleInput, actual []data.BattleInput) {
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, expFirstInput, actual[0])
	assert.Equal(t, expSecondInput, actual[1])
}

func GetPokemonFromId(partyIndex int) *data.BattlePokemon {
	user, _ := utils.ReadJsonFromFile[data.UserSave]("/testfiles/user/test_user.json")
	pokemon := user.ToUser().Party[partyIndex]
	return &data.BattlePokemon{
		Pokemon:      pokemon,
		BattleHP:     pokemon.Stats.HP.Value,
		PokemonFaced: make([]data.BattlePokemon, 0),
	}
}
