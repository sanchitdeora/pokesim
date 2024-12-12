package battle

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

func GetTurnOrder(userInput data.BattleInput, opponentInput data.BattleInput) []data.BattleInput {
	// we will give preference to user whenever equal priority

	if userInput.Type == data.Run || opponentInput.Type == data.Run {
		slog.Error("cannot run in a trainer battle. Need to implement logic here")
		panic("implement logic here")
	} else if userInput.Type == data.Switch || userInput.Type == data.Bag {
		return []data.BattleInput{userInput, opponentInput}
	} else if opponentInput.Type == data.Switch || opponentInput.Type == data.Bag {
		return []data.BattleInput{opponentInput, userInput}
	} else {
		if userInput.Move.Priority != opponentInput.Move.Priority {
			if userInput.Move.Priority > opponentInput.Move.Priority {
				return []data.BattleInput{userInput, opponentInput}
			} else {
				return []data.BattleInput{opponentInput, userInput}
			}
		} else {
			userSpeed := battleStatCalculator(userInput.Selected.Stats.Speed, userInput.Selected.Level)
			targetSpeed := battleStatCalculator(opponentInput.Selected.Stats.Speed, opponentInput.Selected.Level)

			if userSpeed >= targetSpeed {
				return []data.BattleInput{userInput, opponentInput}
			} else {
				return []data.BattleInput{opponentInput, userInput}
			}
		}
	}
}

func battleStatCalculator(stat data.PokemonStat, level int) float64 {
	return (((float64(2*stat.Value) + float64(stat.IV) + float64(stat.EV/4)) * float64(level)) / 100) + 5
}

func battleHPCalculator(HP data.PokemonStat, level int) float64 {
	return (((float64(2*HP.Value) + float64(HP.IV) + float64(HP.EV/4)) * float64(level)) / 100) + 10
}
