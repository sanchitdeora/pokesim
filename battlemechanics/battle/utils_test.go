package battle_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sanchitdeora/PokeSim/battlemechanics/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetTurnOrder_BothSwitch(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Switch,
		Selected: GetPokemonFromPartyIndex(0),
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Switch,
		Selected: GetPokemonFromPartyIndex(1),
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, userAction.ID, opponentAction.ID, order)
}

func TestGetTurnOrder_BothUseBag(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Bag,
		Selected: GetPokemonFromPartyIndex(0),
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Bag,
		Selected: GetPokemonFromPartyIndex(1),
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, userAction.ID, opponentAction.ID, order)
}

func TestGetTurnOrder_BothAttack_UserHighPriorityMove(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(1),
		Move: &data.Moves{
			Priority: -1,
		},
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, userAction.ID, opponentAction.ID, order)
}

func TestGetTurnOrder_BothAttack_OpponentHighPriorityMove(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(0),
		Move: &data.Moves{
			Priority: -1,
		},
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, opponentAction.ID, userAction.ID, order)
}

func TestGetTurnOrder_BothAttack_UserHighSpeed(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentAction := data.BattleAction{
		ID: uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, userAction.ID, opponentAction.ID, order)
}

func TestGetTurnOrder_BothAttack_OpponentHighSpeed(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(0),
		Move: &data.Moves{
			Name:     "Flamethrower",
			Priority: 1,
		},
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(1),
		Move: &data.Moves{
			Name:     "Water Gun",
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, opponentAction.ID, userAction.ID, order)
}

func TestGetTurnOrder_UserAttackOpponentSwitch_UserHighSpeed(t *testing.T) {
	userAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Attack,
		Selected: GetPokemonFromPartyIndex(1),
		Move: &data.Moves{
			Priority: 1,
		},
	}
	opponentAction := data.BattleAction{
		ID:       uuid.NewString(),
		Type:     data.Switch,
		Selected: GetPokemonFromPartyIndex(0),
		Move: &data.Moves{
			Priority: 1,
		},
	}

	order := battle.GetTurnOrder(userAction, opponentAction)

	assertTurnOrder(t, opponentAction.ID, userAction.ID, order)
}

func assertTurnOrder(t *testing.T, expFirstActionID, expSecondActionID string, actual []data.BattleAction) {
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, expFirstActionID, actual[0].ID)
	assert.Equal(t, expSecondActionID, actual[1].ID)
}

func GetPokemonFromPartyIndex(partyIndex int) *data.BattlePokemon {
	user, _ := utils.ReadJsonFromFile[data.UserSave]("/testfiles/user/test_user.json")
	pokemon := user.ToUser().Party[partyIndex]
	return &data.BattlePokemon{
		Pokemon:      pokemon,
		BattleHP:     pokemon.Stats.HP.Value,
		PokemonFaced: make([]string, 0),
	}
}
