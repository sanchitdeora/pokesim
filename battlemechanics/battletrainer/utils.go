package battletrainer

import (
	"log/slog"
	"math"
	"math/rand"

	"github.com/sanchitdeora/PokeSim/data"
)

// switchPokemon takes an active pokemon and a target pokemon to switch to, and a party to pick from
// if the target pokemon is not fainted, it is returned as the new active pokemon and is removed from the party
// if target pokemon is fainted, the active pokemon is not switched
// if the target pokemon is not in the party, the active pokemon is not switched
func switchPokemon(active *data.BattlePokemon, target *data.BattlePokemon, party []*data.BattlePokemon) (*data.BattlePokemon, []*data.BattlePokemon) {
	for i, pokemon := range party {
		if pokemon.PokemonUUID == target.PokemonUUID {

			// only get active pokemon if unfainted pokemon available; else add to the list
			if pokemon.BattleHP > 0 {
				party = append(party[:i], party[i+1:]...)
				party = append(party, active)
				return pokemon, party
			}
		}
	}
	return active, party
}

// healPokemon restores the BattleHP of the target Pokemon using the given item.
// If the Pokemon's current BattleHP is already at its maximum, the function returns immediately.
// Otherwise, it increases the BattleHP by the item's attribute value.
// If the resulting BattleHP exceeds the maximum allowed HP calculated by BattleHPCalculator,
// it is capped at that maximum value.
func healPokemon(targetPokemon *data.BattlePokemon, item *data.Item) bool {
	if item.Category != data.MedicalItems || targetPokemon.BattleHP == targetPokemon.Pokemon.Stats.HP.Value {
		return false
	}
	targetPokemon.BattleHP = min(targetPokemon.BattleHP+int(item.Attribute), targetPokemon.Pokemon.Stats.HP.Value)

	slog.Info("healing pokemon", "pokemon", targetPokemon.Pokemon.Name, "Health:", targetPokemon.BattleHP, "item", item.Description)
	return true
}

func catchPokemon(targetPokemon *data.BattlePokemon, item *data.Item, pokemonCaughtCount int) (pokemonCaught bool, shakeCount int) {
	if item.Category != data.PokeBalls {
		return false, 0
	} 
	catchRate := 1.0

	catchRate = ((3 * float64(targetPokemon.Stats.HP.Value)) - (2 * float64(targetPokemon.BattleHP)))

	catchRate *= float64(targetPokemon.CaptureRate) * item.Attribute
	catchRate = math.Floor(catchRate / (3 * float64(targetPokemon.Stats.HP.Value)))

	// check if catch is successful immediately
	if catchRate * 4096 > 1044480 {
		return true, 0
	}

	// check for critical capture
	critCaptureRate := (catchRate * getCriticalCaptureMultiplier(pokemonCaughtCount)) / 6.0
	if randomGenerator(0, 255) < critCaptureRate {
		return true, 1
	}

	shakeProb := math.Floor(65536 * math.Pow(float64((catchRate*4096)/1044480), 0.1875))
	shakeCount = 1
	for i := 0; i < 4; i++ {
		if randomGenerator(0, 65536) >= shakeProb {
			return false, shakeCount
		}
		shakeCount++
	}

	return true, shakeCount
}

func getCriticalCaptureMultiplier(caught int) float64 {
	if caught > 600 {
		return 2.5
	} else if caught > 450 {
		return 2
	} else if caught > 300 {
		return 1.5
	} else if caught > 150 {
		return 1
	} else if caught > 30 {
		return 0.5
	} else {
		return 0
	}
}

func randomMove(moveset *data.Moveset) *data.Moves {
	moves := []*data.Moves{}
	if moveset.Move1 != nil {
		moves = append(moves, moveset.Move1)
	}
	if moveset.Move2 != nil {
		moves = append(moves, moveset.Move2)
	}
	if moveset.Move3 != nil {
		moves = append(moves, moveset.Move3)
	}
	if moveset.Move4 != nil {
		moves = append(moves, moveset.Move4)
	}

	if len(moves) == 0 {
		return nil
	}

	input := randomGenerator(0, float64(len(moves)))
	return moves[int(input)]
}

func randomGenerator(min float64, max float64) float64 {
	randIndex := rand.Float64()
	randomGenerator := (min + randIndex*(max-min))
	if randomGenerator == max {
		return randomGenerator - 1
	}
	return randomGenerator
}
