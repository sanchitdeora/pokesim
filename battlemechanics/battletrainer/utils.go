package battletrainer

import "github.com/sanchitdeora/PokeSim/data"

func createBattlePokemon(pokemon data.Pokemon) *data.BattlePokemon {
	return &data.BattlePokemon{
		Pokemon: pokemon,
		BattleHP: pokemon.Stats.HP.Value,
		PokemonFaced: make([]data.BattlePokemon, 0),
	}
}

func switchPokemon(active *data.BattlePokemon, target *data.BattlePokemon, party []data.BattlePokemon) (*data.BattlePokemon, []data.BattlePokemon) {
	for i, pokemon := range party {
		if pokemon.ID == target.ID {

			// only get active pokemon if unfainted pokemon available; else add to the list
			if !pokemon.IsFainted {
				party = append(party[:i], party[i+1:]...)
			}
			party = append(party, *active)

			return &pokemon, party
		}
	}
	return active, party
}
