package battle

import (
	"log/slog"
	"math"
	"math/rand"

	"github.com/sanchitdeora/PokeSim/data"
)

func GetTurnOrder(userAction data.BattleAction, opponentAction data.BattleAction) []data.BattleAction {
	// we will give preference to user whenever equal priority

	if userAction.Type == data.Run || opponentAction.Type == data.Run {
		slog.Error("cannot run in a trainer battle. Need to implement logic here")
		panic("implement logic here")
	} else if userAction.Type == data.Switch || userAction.Type == data.Bag {
		return []data.BattleAction{userAction, opponentAction}
	} else if opponentAction.Type == data.Switch || opponentAction.Type == data.Bag {
		return []data.BattleAction{opponentAction, userAction}
	} else {
		if userAction.Move.Priority != opponentAction.Move.Priority {
			if userAction.Move.Priority > opponentAction.Move.Priority {
				return []data.BattleAction{userAction, opponentAction}
			} else {
				return []data.BattleAction{opponentAction, userAction}
			}
		} else {
			if userAction.Selected.Stats.Speed.Value >= opponentAction.Selected.Stats.Speed.Value {
				return []data.BattleAction{userAction, opponentAction}
			} else {
				return []data.BattleAction{opponentAction, userAction}
			}
		}
	}
}

func GetNextUnfaintedPokemonAndCount(party []*data.BattlePokemon) (index, count int) {
	index = -1
	for i, p := range party {
		if p.BattleHP != 0 {
			count ++
			if index < 0 {
				index = i
			}
		}
	}
	return index, count
}

func IsRunSuccessful(pokemon, target *data.BattlePokemon, attempt int) bool {
	if pokemon.Stats.Speed.Value >= target.Stats.Speed.Value {
		return true
	}

	odd := (((float64(pokemon.Stats.Speed.Value) * 32.0) / (float64(target.Stats.Speed.Value) / 4.0)) + 30.0 * float64(attempt)) / 256.0

	return math.Floor(rand.Float64() * 256) < odd 
}